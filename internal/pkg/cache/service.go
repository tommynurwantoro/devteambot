package cache

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNil .
	ErrNil = errors.New("cache: value doesnt exists")
	// ErrCache .
	ErrCache = errors.New("cache error")
)

type Service interface {
	Get(ctx context.Context, key string, value interface{}) error
	Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, value map[string]interface{}) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (bool, error)
	Increment(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string, value int64) (int64, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Ping(ctx context.Context) error
	Close() error
}
