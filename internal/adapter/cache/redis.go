package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

func (r *Redis) Startup() error  { return nil }
func (r *Redis) Shutdown() error { return r.Client.Close() }

// Get .
func (c *Redis) Get(ctx context.Context, key string, value interface{}) error {
	values, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrNil
	} else if err != nil {
		return err
	}
	return json.Unmarshal([]byte(values), &value)
}

// Put .
func (c *Redis) Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Client.Set(ctx, key, string(bytes), expiration).Err()
}

// Expire .
func (c *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

// Delete .
func (c *Redis) Delete(ctx context.Context, keys ...string) (int64, error) {
	return c.Client.Del(ctx, keys...).Result()
}

func (c *Redis) Exists(ctx context.Context, keys ...string) (bool, error) {
	exist, err := c.Client.Exists(ctx, keys...).Result()
	if err != nil {
		return false, err
	}

	if exist > 0 {
		return true, nil
	}

	return false, nil
}

func (c *Redis) Increment(ctx context.Context, key string, value int64) (int64, error) {
	return c.Client.IncrBy(ctx, key, value).Result()
}

func (c *Redis) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return c.Client.DecrBy(ctx, key, value).Result()
}

func (c *Redis) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := c.Client.Keys(ctx, pattern).Result()
	return keys, err
}

func (c *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.Client.TTL(ctx, key).Result()
}

// Ping .
func (c *Redis) Ping(ctx context.Context) error {
	_, err := c.Client.Ping(ctx).Result()
	return err
}

// Close .
func (c *Redis) Close() error {
	return c.Client.Close()
}
