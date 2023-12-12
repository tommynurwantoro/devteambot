package scheduler

import (
	"context"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/repository/gorm"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/adapter/resty"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type SholatScheduler struct {
	Scheduler         *Scheduler         `inject:"scheduler"`
	Discord           *discord.App       `inject:"discord"`
	MyQuranAPI        *resty.MyQuran     `inject:"myQuranAPI"`
	RedisKey          redis.RedisKey     `inject:"redisKey"`
	SettingKey        gorm.SettingKey    `inject:"settingKey"`
	Cache             cache.Service      `inject:"cache"`
	SettingRepository setting.Repository `inject:"settingRepository"`
}

func (s *SholatScheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	s.Scheduler.Every(1).Day().At("03:00").Do(func() {
		logger.Info("Get Sholat Schedule")
		s.GetSholatSchedule(context.Background())
	})

	s.Scheduler.Every(1).Minute().Do(func() {
		now := time.Now().In(loc)
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			return
		}
		s.SendReminderSholat(context.Background())
	})

	return nil
}

func (s *SholatScheduler) Shutdown() error { return nil }

func (s *SholatScheduler) GetSholatSchedule(ctx context.Context) {
	response := new(resty.GetSholatResponse)
	req := s.MyQuranAPI.Client.R().SetContext(ctx).
		ForceContentType("application/json").
		SetResult(response)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	_, err := req.Get(fmt.Sprintf("/sholat/jadwal/1505/%d/%d/%d", now.Year(), int(now.Month()), now.Day()))
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %s", err.Error()), err)
		return
	}

	s.Cache.Put(ctx, s.RedisKey.DailySholatSchedule(), response, 24*time.Hour)
}

func (s *SholatScheduler) SendReminderSholat(ctx context.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	sholatSchedule := resty.GetSholatResponse{}
	s.Cache.Get(ctx, s.RedisKey.DailySholatSchedule(), &sholatSchedule)

	timeNow := now.Format("15:04")
	if timeNow != sholatSchedule.Data.Jadwal.Dzuhur && timeNow != sholatSchedule.Data.Jadwal.Ashar {
		return
	}

	settings, err := s.SettingRepository.GetAllByKey(ctx, s.SettingKey.ReminderSholat())
	if err != nil {
		return
	}

	if len(settings) == 0 {
		return
	}

	for _, set := range settings {
		var val, channelID, roleID string
		err = json.Unmarshal([]byte(set.Value), &val)
		if err != nil {
			logger.Error("Error: "+err.Error(), err)
			return
		}

		split := strings.Split(val, "|")
		channelID = split[0]
		roleID = split[1]

		logger.Info("Send reminder sholat")
		s.Discord.Bot.ChannelMessageSend(channelID, fmt.Sprintf("Udah adzan, yuk sholat dulu <@&%s> !", roleID))
	}
}
