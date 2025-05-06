package cache

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

// ProductIndexName es el nombre del índice de productos en RediSearch
const ProductIndexName = "idx:producto"

// RediSearchClient define las operaciones de búsqueda sobre RediSearch
type RediSearchClient interface {
	CreateIndex(ctx context.Context, name string, schema map[string]string) error
	AddDocument(ctx context.Context, index, id string, fields map[string]interface{}) error
	Search(ctx context.Context, index, query string, offset, limit int) ([]map[string]string, int64, error)
	DeleteDocument(ctx context.Context, index, id string) error
}

type clientImpl struct {
	client      *redisearch.Client
	pool        *redis.Pool // add pool for direct connections
	isAvailable bool        // Track if RediSearch is available
}

// NewRediSearchClient crea un cliente de RediSearch a partir de la URL
func NewRediSearchClient(redisURL string) (RediSearchClient, error) {
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error al parsear URL de Redis: %w", err)
	}

	host := parsedURL.Host
	if host == "" {
		host = "localhost:6379"
	}

	password, _ := parsedURL.User.Password()

	db := 0
	if path := parsedURL.Path; path != "" {
		pathParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
		if len(pathParts) > 0 && pathParts[0] != "" {
			if dbNum, err := strconv.Atoi(pathParts[0]); err == nil {
				db = dbNum
			}
		}
	}

	dialOpts := []redis.DialOption{
		redis.DialConnectTimeout(5 * time.Second),
		redis.DialReadTimeout(5 * time.Second),
		redis.DialWriteTimeout(5 * time.Second),
	}

	if password != "" {
		dialOpts = append(dialOpts, redis.DialPassword(password))
	}

	if db > 0 {
		dialOpts = append(dialOpts, redis.DialDatabase(db))
	}

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host, dialOpts...)
		},
		MaxIdle:     10,
		IdleTimeout: 5 * time.Minute,
		MaxActive:   100,
	}

	client := redisearch.NewClientFromPool(pool, ProductIndexName)
	impl := &clientImpl{
		client:      client,
		pool:        pool, // store pool
		isAvailable: false,
	}

	// Test if RediSearch is available
	if err := impl.testRediSearch(); err != nil {
		return nil, fmt.Errorf("RediSearch not available: %w", err)
	}
	impl.isAvailable = true
	return impl, nil
}

func (c *clientImpl) testRediSearch() error {
	conn := c.pool.Get()
	defer conn.Close()

	// 1) Intentar detectar módulo vía MODULE LIST
	modules, err := redis.Values(conn.Do("MODULE", "LIST"))
	if err == nil {
		for _, m := range modules {
			if mi, ok := m.([]interface{}); ok {
				for i := 0; i < len(mi)-1; i += 2 {
					if key, _ := redis.String(mi[i], nil); key == "name" {
						if name, _ := redis.String(mi[i+1], nil); strings.Contains(strings.ToLower(name), "search") {
							return nil
						}
					}
				}
			}
		}
		return fmt.Errorf("RediSearch module not found in MODULE LIST")
	}

	// 2) Fallback: sonda con FT.SEARCH en índice inexistente
	_, err = conn.Do("FT.SEARCH", "redisearch_test", "*", "LIMIT", 0, 0)
	if err != nil {
		msg := err.Error()
		if strings.HasPrefix(msg, "ERR unknown command") {
			return fmt.Errorf("RediSearch module not loaded")
		}
		// cualquier otro error indica que el módulo está presente
	}
	return nil
}

func (c *clientImpl) CreateIndex(ctx context.Context, name string, schema map[string]string) error {
	if !c.isAvailable {
		log.Printf("WARNING: Skipping index creation - RediSearch not available")
		return nil // Skip index creation if RediSearch isn't available
	}

	// 1) chequeo explícito de existencia
	conn := c.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("FT.INFO", name); err != nil {
		if strings.HasPrefix(err.Error(), "ERR unknown index name") {
			// índice no existe → no hay nada que dropear
		} else {
			return fmt.Errorf("error checking index existence: %w", err)
		}
	} else {
		// 2) índice existe → dropear y propagar fallo si lo hay
		if err := c.client.DropIndex(true); err != nil {
			return fmt.Errorf("failed to drop existing index: %w", err)
		}
	}

	// Create schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions)

	for field, typ := range schema {
		switch typ {
		case "TEXT":
			sc = sc.AddField(redisearch.NewTextFieldOptions(field, redisearch.TextFieldOptions{Sortable: true}))
		case "NUMERIC":
			sc = sc.AddField(redisearch.NewNumericFieldOptions(field, redisearch.NumericFieldOptions{Sortable: true}))
		case "TAG":
			sc = sc.AddField(redisearch.NewTagFieldOptions(field, redisearch.TagFieldOptions{Separator: ',', Sortable: true}))
		default:
			return fmt.Errorf("unsupported field type: %s", typ)
		}
	}

	// Try to create index but handle case where RediSearch is not available
	if err := c.client.CreateIndex(sc); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

func (c *clientImpl) AddDocument(ctx context.Context, _index, id string, fields map[string]interface{}) error {
	if !c.isAvailable {
		log.Printf("WARNING: Skipping AddDocument - RediSearch not available")
		return nil
	}
	doc := redisearch.NewDocument(id, 1.0)
	for k, v := range fields {
		doc.Set(k, v)
	}
	return c.client.IndexOptions(redisearch.IndexingOptions{Replace: true}, doc)
}

func (c *clientImpl) Search(ctx context.Context, _index, query string, offset, limit int) ([]map[string]string, int64, error) {
	if !c.isAvailable {
		log.Printf("WARNING: Skipping Search - RediSearch not available")
		return nil, 0, nil
	}
	q := redisearch.NewQuery(query).Limit(offset, limit)
	docs, total, err := c.client.Search(q)
	if err != nil {
		return nil, 0, err
	}
	results := make([]map[string]string, 0, len(docs))
	for _, d := range docs {
		m := map[string]string{"id": d.Id}
		for k, v := range d.Properties {
			m[k] = fmt.Sprint(v)
		}
		results = append(results, m)
	}
	return results, int64(total), nil
}

func (c *clientImpl) DeleteDocument(ctx context.Context, _index, id string) error {
	if !c.isAvailable {
		log.Printf("WARNING: Skipping DeleteDocument - RediSearch not available")
		return nil
	}
	return c.client.Delete(id, true)
}
