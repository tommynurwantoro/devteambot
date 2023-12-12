package scheduler

import (
	"context"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/repository/gorm"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type PresensiScheduler struct {
	Scheduler         *Scheduler         `inject:"scheduler"`
	Discord           *discord.App       `inject:"discord"`
	SettingKey        gorm.SettingKey    `inject:"settingKey"`
	SettingRepository setting.Repository `inject:"settingRepository"`
}

func (s *PresensiScheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	s.Scheduler.Every(1).Day().At("07:55").Do(func() {
		now := time.Now().In(loc)
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			return
		}

		s.SendReminderPresensi(context.Background())
	})

	s.Scheduler.Every(1).Day().At("17:05").Do(func() {
		now := time.Now().In(loc)
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			return
		}

		s.SendReminderPresensi(context.Background())
	})

	return nil
}

func (s *PresensiScheduler) Shutdown() error { return nil }

func (s *PresensiScheduler) SendReminderPresensi(ctx context.Context) {
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
		s.Discord.Bot.ChannelMessageSend(channelID, fmt.Sprintf("Jangan lupa melakukan presensi <@&%s> !", roleID))
	}
}
