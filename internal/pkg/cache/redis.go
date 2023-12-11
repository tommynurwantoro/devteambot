package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

// New redis
func New(conf RedisConfig) *Redis {
	option := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password: conf.Password,
	}

	client := redis.NewClient(option)

	return &Redis{Client: client}
}

// Get data for the given key
func (c *Redis) Get(ctx context.Context, key string, value interface{}) error {
	values, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrNil
	} else if err != nil {
		return err
	}
	return json.Unmarshal([]byte(values), &value)
}

// Put data to redis based on key
func (c *Redis) Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Client.Set(ctx, key, string(bytes), expiration).Err()
}

func (c *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	values, err := c.Client.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return values, ErrNil
	} else if err != nil {
		return values, err
	}

	if len(values) == 0 {
		return values, ErrNil
	}

	return values, nil
}

func (c *Redis) HSet(ctx context.Context, key string, value map[string]interface{}) error {
	return c.Client.HSet(ctx, key, value).Err()
}

// Set expiration for the given key
func (c *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

// Delete data based on key
func (c *Redis) Delete(ctx context.Context, keys ...string) (int64, error) {
	return c.Client.Del(ctx, keys...).Result()
}

// Check if key exist in Redis
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

// Increase value for the given key
func (c *Redis) Increment(ctx context.Context, key string, value int64) (int64, error) {
	return c.Client.IncrBy(ctx, key, value).Result()
}

// Decrease value for the given key
func (c *Redis) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return c.Client.DecrBy(ctx, key, value).Result()
}

// Get available key for the given pattern
func (c *Redis) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := c.Client.Keys(ctx, pattern).Result()
	return keys, err
}

// Check expiration for the given key
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
