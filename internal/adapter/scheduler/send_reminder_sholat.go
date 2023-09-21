package scheduler

import (
	"context"
	"devteambot/internal/adapter/resty"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (s *Scheduler) SendReminderSholat(ctx context.Context) {
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
		s.App.Bot.ChannelMessageSend(channelID, fmt.Sprintf("Udah adzan, yuk sholat dulu <@&%s> !", roleID))
	}
}
