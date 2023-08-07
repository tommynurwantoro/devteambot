package scheduler

import (
	"context"
	"devteambot/internal/adapter/resty"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"time"
)

func (s *Scheduler) SendReminderSholat(ctx context.Context) {
	now := time.Now()
	sholatSchedule := resty.GetSholatResponse{}
	s.Cache.Get(ctx, s.RedisKey.DailySholatSchedule(), &sholatSchedule)

	timeNow := now.Format("15:04")
	if timeNow != sholatSchedule.Data.Jadwal.Dzuhur && timeNow != sholatSchedule.Data.Jadwal.Ashar {
		return
	}

	settings, err := s.SettingRepository.GetAllByKey(ctx, s.SettingKey.ReminderSholatChannel())
	if err != nil {
		return
	}

	if len(settings) == 0 {
		return
	}

	for _, set := range settings {
		var channelID string
		err = json.Unmarshal([]byte(set.Value), &channelID)
		if err != nil {
			logger.Error("Error: "+err.Error(), err)
			return
		}

		logger.Info("Send reminder sholat")
		s.App.Bot.ChannelMessageSend(channelID, "Udah adzan, yuk sholat dulu @everyone !")
	}
}
