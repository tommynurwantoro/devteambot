package bootstrap

import (
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/google"
	"devteambot/internal/adapter/n8n"
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

func RegisterAI() {
	appContainer.RegisterService("googleai", new(google.AI))
	appContainer.RegisterService("n8n", new(n8n.N8N))
}

func RegisterRest() {
	appContainer.RegisterService("rest", new(rest.Fiber))
}

func RegisterDiscord() {
	appContainer.RegisterService("discord", new(discord.App))
}

func RegisterResty() {
	appContainer.RegisterService("jadwalSholatOrgAPI", new(resty.JadwalSholatOrg))
}
