package scheduler

import (
	"context"
	"devteambot/internal/pkg/logger"
	"encoding/json"
)

func (s *Scheduler) SendReminderPresensi(ctx context.Context) {
	settings, err := s.SettingRepository.GetAllByKey(ctx, s.SettingKey.ReminderPresensiChannel())
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

		logger.Info("Send reminder presensi")
		s.App.Bot.ChannelMessageSend(channelID, "Jangan lupa melakukan presensi @everyone !")
	}
}
