package events

import (
	"devteambot/internal/adapter/discord"
	"devteambot/internal/pkg/cache"
)

type Event struct {
	Discord      *discord.App    `inject:"discord"`
	Cache        cache.Service   `inject:"cache"`
	Admins       map[string]bool `inject:"admins"`
	RandomCoinAt int64
}

func (e *Event) Startup() error { return nil }

func (e *Event) Shutdown() error { return nil }
