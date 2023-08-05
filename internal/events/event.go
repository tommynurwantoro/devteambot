package events

import (
	"devteambot/config"
	"devteambot/internal/adapter/cache"
	"devteambot/internal/constant"

	"github.com/bwmarrin/discordgo"
)

type Event struct {
	Session      *discordgo.Session `inject:"botSession"`
	Cache        cache.Cache        `inject:"cache"`
	Conf         config.Discord     `inject:"discordConfig"`
	Key          constant.RedisKey  `inject:"redisKey"`
	Color        constant.Color     `inject:"color"`
	Admins       map[string]bool    `inject:"admins"`
	RandomCoinAt int64
}

func (e *Event) Startup() error { return nil }

func (e *Event) Shutdown() error { return nil }
