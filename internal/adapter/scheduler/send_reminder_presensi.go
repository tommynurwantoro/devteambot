package scheduler

import (
	"context"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
)

func (s *Scheduler) SendReminderPresensi(ctx context.Context) {
	settings, err := s.SettingRepository.GetAllByKey(ctx, s.SettingKey.ReminderPresensi())
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

		logger.Info("Send reminder presensi")
		s.App.Bot.ChannelMessageSend(channelID, fmt.Sprintf("Jangan lupa melakukan presensi <@&%s> !", roleID))
	}
}
