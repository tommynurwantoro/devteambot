package scheduler

import (
	"time"

	"devteambot/config"
	"devteambot/internal/adapter/cache"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/constant"
	"devteambot/internal/domain/setting"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	Conf       config.Discord      `inject:"discordConfig"`
	Cache      cache.Cache         `inject:"cache"`
	App        *discord.App        `inject:"discordApp"`
	SettingKey constant.SettingKey `inject:"settingKey"`
	RedisKey   constant.RedisKey   `inject:"redisKey"`
	Color      constant.Color      `inject:"color"`

	SettingRepository setting.Repository `inject:"settingRepository"`
}

func (s *Scheduler) Startup() error {
	// ctx := context.Background()
	scheduler := gocron.NewScheduler(time.UTC)

	// Archive mod report everyday
	// scheduler.Every(30).Seconds().Do(func() {
	scheduler.Every(1).Day().At("00:00").Do(func() {
		// s.SendModReport(ctx)
	})

	scheduler.StartAsync()

	return nil
}

func (s *Scheduler) Shutdown() error { return nil }

func (s *Scheduler) contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}

	return false
}
