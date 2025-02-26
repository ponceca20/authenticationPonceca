package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultExpiration   = 24 * time.Hour
	DefaultPoolSize     = 500
	DefaultMinIdleConns = 100
	DefaultMaxRetries   = 3
	DefaultRetryBackoff = 50 * time.Millisecond
)

type RedisCache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

type redisCacheImpl struct {
	client *redis.Client
	mu     sync.RWMutex
}

// Configuración optimizada para múltiples empresas
func NewRedisCache(redisURL string) (RedisCache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error al analizar la URL de Redis: %w", err)
	}

	// Configuración optimizada para múltiples empresas
	opt.PoolSize = DefaultPoolSize
	opt.MinIdleConns = DefaultMinIdleConns
	opt.MaxRetries = DefaultMaxRetries
	opt.MaxRetryBackoff = DefaultRetryBackoff
	opt.ReadTimeout = 2 * time.Second  // Aumentado para mejor tolerancia
	opt.WriteTimeout = 2 * time.Second // Aumentado para mejor tolerancia
	opt.PoolTimeout = 4 * time.Second  // Aumentado para mejor tolerancia
	opt.ConnMaxIdleTime = 10 * time.Minute
	opt.DialTimeout = 5 * time.Second // Añadido timeout de conexión

	client := redis.NewClient(opt)

	// Verificar conexión con retry
	var lastErr error
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := client.Ping(ctx).Err()
		cancel()

		if err == nil {
			fmt.Println("✅ Conexión exitosa a Redis")
			return &redisCacheImpl{client: client}, nil
		}

		lastErr = err
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return nil, fmt.Errorf("no se pudo conectar a Redis después de 3 intentos: %w", lastErr)
}

// Implementación optimizada de Get con timeout y retry
func (r *redisCacheImpl) Get(ctx context.Context, key string, dest interface{}) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Usar un contexto con timeout más largo
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var lastErr error
	for i := 0; i < DefaultMaxRetries; i++ {
		data, err := r.client.Get(ctxTimeout, key).Bytes()
		if err == redis.Nil {
			return nil
		}
		if err != nil {
			lastErr = err
			time.Sleep(DefaultRetryBackoff)
			continue
		}

		if err := json.Unmarshal(data, dest); err != nil {
			return fmt.Errorf("error al deserializar datos: %w", err)
		}
		return nil
	}

	return fmt.Errorf("error al obtener datos de Redis después de %d intentos: %w", DefaultMaxRetries, lastErr)
}

// Implementación optimizada de Set con retry
func (r *redisCacheImpl) Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Usar un contexto con timeout más largo
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error al serializar datos: %w", err)
	}

	exp := DefaultExpiration
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	var lastErr error
	for i := 0; i < DefaultMaxRetries; i++ {
		if err := r.client.Set(ctxTimeout, key, data, exp).Err(); err != nil {
			lastErr = err
			time.Sleep(DefaultRetryBackoff)
			continue
		}
		return nil
	}

	return fmt.Errorf("error al guardar en Redis después de %d intentos: %w", DefaultMaxRetries, lastErr)
}

func (r *redisCacheImpl) Delete(ctx context.Context, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var lastErr error
	for i := 0; i < DefaultMaxRetries; i++ {
		if err := r.client.Del(ctxTimeout, key).Err(); err != nil {
			lastErr = err
			time.Sleep(DefaultRetryBackoff)
			continue
		}
		return nil
	}

	return fmt.Errorf("error al eliminar de Redis después de %d intentos: %w", DefaultMaxRetries, lastErr)
}

func (r *redisCacheImpl) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.client.Close()
}

// NewRedisCacheFromEnv creates a new RedisCache instance using environment variables
func NewRedisCacheFromEnv() (RedisCache, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, fmt.Errorf("REDIS_URL no está configurada en las variables de entorno")
	}

	cache, err := NewRedisCache(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error al crear la caché de Redis: %w", err)
	}

	return cache, nil
}
