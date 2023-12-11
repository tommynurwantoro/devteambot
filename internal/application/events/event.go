package events

import (
	"devteambot/internal/adapter/discord"
	"devteambot/internal/constant"
	"devteambot/internal/pkg/cache"
)

type Event struct {
	Discord      *discord.App      `inject:"discord"`
	Cache        cache.Service     `inject:"cache"`
	Key          constant.RedisKey `inject:"redisKey"`
	Color        constant.Color    `inject:"color"`
	Admins       map[string]bool   `inject:"admins"`
	RandomCoinAt int64
}

func (e *Event) Startup() error { return nil }

func (e *Event) Shutdown() error { return nil }
