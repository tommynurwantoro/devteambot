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
	API               *resty.JadwalSholatOrg `inject:"jadwalSholatOrgAPI"`
	Cache             cache.Service          `inject:"cache"`
	App               *discord.App           `inject:"discord"`
	SettingRepository setting.Repository     `inject:"settingRepository"`
	RedisKey          redis.RedisKey         `inject:"redisKey"`
	SettingKey        gorm.SettingKey        `inject:"settingKey"`
}

func (s *SholatService) Startup() error { return nil }

func (s *SholatService) Shutdown() error { return nil }

func (s *SholatService) GetSholatSchedule(ctx context.Context) error {
	logger.Info("Get Sholat Schedule")
	req := s.API.Client.R().SetContext(ctx).
		ForceContentType("application/json")
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	resp, err := req.Get(fmt.Sprintf("yogyakarta/%d/%d.json", now.Year(), int(now.Month())))
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %s", err.Error()), err)
		return err
	}

	response := make([]resty.GetJadwalSholatResponse, 0)
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %s", err.Error()), err)
		return err
	}

	todaySchedule := response[now.Day()-1]
	s.Cache.Put(ctx, s.RedisKey.DailySholatSchedule(), todaySchedule, 24*time.Hour)
	return nil
}

func (s *SholatService) SendReminderSholat(ctx context.Context) error {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	sholatSchedule := make([]resty.GetJadwalSholatResponse, 0)
	s.Cache.Get(ctx, s.RedisKey.DailySholatSchedule(), &sholatSchedule)

	if len(sholatSchedule) == 0 {
		return nil
	}

	timeNow := now.Format("15:04")
	todaySchedule := sholatSchedule[now.Day()-1]
	if timeNow != todaySchedule.Dzuhur && timeNow != todaySchedule.Ashar {
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
