package service

import (
	"context"
	"devteambot/internal/adapter/discord"
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
}

func (s *SholatService) Startup() error { return nil }

func (s *SholatService) Shutdown() error { return nil }

func (s *SholatService) GetTodaySchedule(ctx context.Context) error {
	req := s.API.Client.R().SetContext(ctx).
		ForceContentType("application/json")
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	logger.Info(fmt.Sprintf("Sholat: Get Today Schedule is running at %s", time.Now().In(loc).Format("2006-01-02 15:04:05")))

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

func (s *SholatService) SendReminder(ctx context.Context) error {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	logger.Info(fmt.Sprintf("Sholat: Send Reminder is running at %s", time.Now().In(loc).Format("2006-01-02 15:04:05")))

	var todaySchedule resty.GetJadwalSholatResponse
	s.Cache.Get(ctx, s.RedisKey.DailySholatSchedule(), &todaySchedule)

	timeNow := now.Format("15:04")
	if timeNow != todaySchedule.Dzuhur && timeNow != todaySchedule.Ashar {
		return nil
	}

	settings, err := s.SettingRepository.GetAllByKey(ctx, setting.REMINDER_SHOLAT)
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
