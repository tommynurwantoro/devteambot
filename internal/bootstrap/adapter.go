package bootstrap

import (
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/repository"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/adapter/rest"
	"devteambot/internal/adapter/resty"
)

func RegisterDatabase() {
	appContainer.RegisterService("database", new(repository.Gorm))
}

func RegisterCache() {
	appContainer.RegisterService("cache", new(repository.Cache))
	appContainer.RegisterService("redisKey", redis.NewRedisKey())
}

func RegisterRest() {
	appContainer.RegisterService("rest", new(rest.Fiber))
}

func RegisterDiscord() {
	appContainer.RegisterService("discord", new(discord.App))
}

func RegisterResty() {
	appContainer.RegisterService("myQuranAPI", new(resty.MyQuran))
}
