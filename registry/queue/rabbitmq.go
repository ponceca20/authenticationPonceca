package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DefaultReconnectDelay = 5 * time.Second
	DefaultResendDelay    = 1 * time.Second
	DefaultTimeout        = 30 * time.Second
	DefaultPrefetchCount  = 1
	DefaultExchangeType   = "topic"
)

// RabbitMQ interface defines the methods for interacting with RabbitMQ
type RabbitMQ interface {
	Publish(ctx context.Context, exchange, routingKey string, message interface{}) error
	Subscribe(queueName, routingKey string, handler func([]byte) error) error
	Close() error
}

type rabbitMQImpl struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mu      sync.RWMutex
	config  RabbitConfig
}

type RabbitConfig struct {
	URL            string
	ReconnectDelay time.Duration
	ResendDelay    time.Duration
	PrefetchCount  int
	ExchangeName   string
	ExchangeType   string
	ConnectionName string
}

// NewRabbitMQFromEnv creates a new RabbitMQ instance using environment variables
func NewRabbitMQFromEnv() (RabbitMQ, error) {
	cfg := RabbitConfig{
		URL:            os.Getenv("RABBITMQ_URL"),
		ExchangeName:   os.Getenv("RABBITMQ_EXCHANGE"),
		ExchangeType:   DefaultExchangeType,
		ReconnectDelay: DefaultReconnectDelay,
		ResendDelay:    DefaultResendDelay,
		PrefetchCount:  DefaultPrefetchCount,
		ConnectionName: fmt.Sprintf("app-%s", uuid.New().String()[:8]),
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}
	if cfg.ExchangeName == "" {
		return nil, fmt.Errorf("RABBITMQ_EXCHANGE is required")
	}

	return NewRabbitMQ(cfg)
}

// NewRabbitMQ creates a new RabbitMQ instance with the provided configuration
func NewRabbitMQ(cfg RabbitConfig) (RabbitMQ, error) {
	if cfg.ReconnectDelay == 0 {
		cfg.ReconnectDelay = DefaultReconnectDelay
	}
	if cfg.ResendDelay == 0 {
		cfg.ResendDelay = DefaultResendDelay
	}
	if cfg.PrefetchCount == 0 {
		cfg.PrefetchCount = DefaultPrefetchCount
	}

	rmq := &rabbitMQImpl{
		config: cfg,
	}

	if err := rmq.connect(); err != nil {
		return nil, err
	}

	return rmq, nil
}

func (r *rabbitMQImpl) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	config := amqp.Config{
		Properties: amqp.Table{
			"connection_name": r.config.ConnectionName,
		},
		Dial: amqp.DefaultDial(DefaultTimeout),
	}

	conn, err := amqp.DialConfig(r.config.URL, config)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS for better load distribution
	if err := ch.Qos(
		r.config.PrefetchCount, // prefetch count
		0,                      // prefetch size
		false,                  // global
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	if err := ch.ExchangeDeclare(
		r.config.ExchangeName,
		r.config.ExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	r.conn = conn
	r.channel = ch

	go r.handleReconnect()
	return nil
}

func (r *rabbitMQImpl) Publish(ctx context.Context, exchange, routingKey string, message interface{}) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	return r.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Headers: amqp.Table{
				"message_id": uuid.New().String(),
			},
		},
	)
}

func (r *rabbitMQImpl) Subscribe(queueName, routingKey string, handler func([]byte) error) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	q, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := r.channel.QueueBind(
		q.Name,
		routingKey,
		r.config.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				// En caso de error, rechazamos el mensaje y lo reencolamos
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		}
	}()

	return nil
}

func (r *rabbitMQImpl) handleReconnect() {
	for {
		reason, ok := <-r.conn.NotifyClose(make(chan *amqp.Error))
		if !ok {
			// La conexión fue cerrada intencionalmente
			return
		}
		fmt.Printf("connection closed, reason: %v\n", reason)

		// Intento de reconexión
		for {
			time.Sleep(r.config.ReconnectDelay)
			if err := r.connect(); err == nil {
				break
			}
		}
	}
}

func (r *rabbitMQImpl) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var errors []error

	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close channel: %w", err))
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple errors occurred during close: %v", errors)
	}

	return nil
}
