package service

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

type SholatService struct {
	MyQuranAPI        *resty.MyQuran     `inject:"myQuranAPI"`
	Cache             cache.Service      `inject:"cache"`
	App               *discord.App       `inject:"discord"`
	SettingRepository setting.Repository `inject:"settingRepository"`
	RedisKey          redis.RedisKey     `inject:"redisKey"`
	SettingKey        gorm.SettingKey    `inject:"settingKey"`
}

func (s *SholatService) GetSholatSchedule(ctx context.Context) error {
	response := new(resty.GetSholatResponse)
	req := s.MyQuranAPI.Client.R().SetContext(ctx).
		ForceContentType("application/json").
		SetResult(response)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	_, err := req.Get(fmt.Sprintf("/sholat/jadwal/1505/%d/%d/%d", now.Year(), int(now.Month()), now.Day()))
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %s", err.Error()), err)
		return err
	}

	s.Cache.Put(ctx, s.RedisKey.DailySholatSchedule(), response, 24*time.Hour)
	return nil
}

func (s *SholatService) SendReminderSholat(ctx context.Context) error {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	sholatSchedule := resty.GetSholatResponse{}
	s.Cache.Get(ctx, s.RedisKey.DailySholatSchedule(), &sholatSchedule)

	timeNow := now.Format("15:04")
	if timeNow != sholatSchedule.Data.Jadwal.Dzuhur && timeNow != sholatSchedule.Data.Jadwal.Ashar {
		return nil
	}

	settings, err := s.SettingRepository.GetAllByKey(ctx, s.SettingKey.ReminderSholat())
	if err != nil {
		return err
	}

	if len(settings) == 0 {
		return nil
	}

	for _, set := range settings {
		var val, channelID, roleID string
		err = json.Unmarshal([]byte(set.Value), &val)
		if err != nil {
			logger.Error("Error: "+err.Error(), err)
			return err
		}

		split := strings.Split(val, "|")
		channelID = split[0]
		roleID = split[1]

		logger.Info("Send reminder sholat")
		s.App.Bot.ChannelMessageSend(channelID, fmt.Sprintf("Udah adzan, yuk sholat dulu <@&%s> !", roleID))
	}

	return nil
}
