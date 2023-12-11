package repository

import (
	"context"
	"devteambot/config"
	"devteambot/internal/pkg/cache"
)

type Cache struct {
	*cache.Redis
	Conf *config.Config `inject:"config"`
}

func (c *Cache) Startup() error {
	redisConfig := cache.RedisConfig{
		Address:  c.Conf.Redis.Address,
		Port:     c.Conf.Redis.Port,
		Password: c.Conf.Redis.Password,
	}

	c.Redis = cache.New(redisConfig)

	return c.Ping(context.Background())
}

func (c *Cache) Shutdown() error {
	return c.Client.Close()
}
