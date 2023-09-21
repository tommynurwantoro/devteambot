package scheduler

import (
	"context"
	"time"

	"devteambot/config"
	"devteambot/internal/adapter/cache"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/resty"
	"devteambot/internal/constant"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"

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
	MyQuranAPI        *resty.MyQuran     `inject:"myQuranAPI"`
}

func (s *Scheduler) Startup() error {
	ctx := context.Background()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := gocron.NewScheduler(loc)

	// scheduler.Every(30).Seconds().Do(func() {
	scheduler.Every(1).Day().At("03:00").Do(func() {
		logger.Info("Get Sholat Schedule")
		s.GetSholatSchedule(ctx)
	})

	scheduler.Every(1).Day().Monday().Tuesday().Wednesday().Thursday().Friday().At("07:55").Do(func() {
		s.SendReminderPresensi(ctx)
	})

	scheduler.Every(1).Day().Monday().Tuesday().Wednesday().Thursday().Friday().At("17:05").Do(func() {
		s.SendReminderPresensi(ctx)
	})

	scheduler.Every(1).Minute().Monday().Tuesday().Wednesday().Thursday().Friday().Do(func() {
		s.SendReminderSholat(ctx)
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
