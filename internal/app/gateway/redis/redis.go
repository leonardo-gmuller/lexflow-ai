package redis

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/config"
	"github.com/redis/go-redis/v9"
	redisv9 "github.com/redis/go-redis/v9"
)

// Interface que a aplicação usa (para facilitar teste / mock)
type ClientInterface interface {
	Ping(ctx context.Context) error

	// List operations
	RPush(ctx context.Context, key string, values ...string) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string, objByRef any) error
	Set(ctx context.Context, key string, obj any, ttl time.Duration) error
	SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error)
	// Get/Set HTTP Response
	GetResp(ctx context.Context, key string) (*http.Response, error)
	SetResp(ctx context.Context, key string, resp *http.Response, ttl time.Duration) error

	// Sorted set operations (para scheduler / jobs)
	ZAdd(ctx context.Context, key string, members ...redisv9.Z) error
	ZRangeByScore(ctx context.Context, key, min, max string, offset, count int64) ([]string, error)
	ZRem(ctx context.Context, key string, members ...string) error

	// Opcional: TTL, se quiser usar depois
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

type Client struct {
	inner *redisv9.Client
}

const errCacheKeyDoesNotExist = redisv9.Nil

// NewClient cria um redis client configurável.
// Exemplo de uso no main:
//
//	redisClient := redis.NewClient("localhost:6379", "", 0)
func NewClient(cfg config.RedisConfig) (*Client, error) {
	addr := net.JoinHostPort(cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping failed addr=%q: %w", addr, err)
	}

	return &Client{inner: rdb}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return c.inner.Ping(ctx2).Err()
}

// ============ List ============

func (c *Client) RPush(ctx context.Context, key string, values ...string) error {
	// RPush aceita ...interface{}, mas aqui tipamos como []string pra ficar mais seguro
	vals := make([]interface{}, len(values))
	for i, v := range values {
		vals[i] = v
	}
	return c.inner.RPush(ctx, key, vals...).Err()
}

func (c *Client) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.inner.LRange(ctx, key, start, stop).Result()
}

// ============ Sorted Set (ZSET) ============

func (c *Client) ZAdd(ctx context.Context, key string, members ...redisv9.Z) error {
	return c.inner.ZAdd(ctx, key, members...).Err()
}

func (c *Client) ZRangeByScore(ctx context.Context, key, min, max string, offset, count int64) ([]string, error) {
	return c.inner.ZRangeByScore(ctx, key, &redisv9.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (c *Client) ZRem(ctx context.Context, key string, members ...string) error {
	vals := make([]interface{}, len(members))
	for i, v := range members {
		vals[i] = v
	}
	return c.inner.ZRem(ctx, key, vals...).Err()
}

// ============ TTL opcional ============

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.inner.Expire(ctx, key, expiration).Err()
}

func (c *Client) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return c.inner.SetNX(ctx, key, value, ttl).Result()
}
