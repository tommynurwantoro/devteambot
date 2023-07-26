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

func (e *Event) Startup() error {
	e.SetRandomCoinAt()

	e.Session.AddHandler(e.Bingo)
	return nil
}

func (e *Event) Shutdown() error { return nil }

func (e *Event) SetRandomCoinAt() {
	// rand.Seed(time.Now().UnixNano())
	// e.RandomCoinAt = rand.Int63n(e.Conf.Gamification.RandomCoin.Every) + 5

	// logger.Info(fmt.Sprintf("SetRandomCoin: %d", e.RandomCoinAt))
}
