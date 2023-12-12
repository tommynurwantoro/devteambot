package service

import (
	"context"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/repository/gorm"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
)

type PresensiService struct {
	Discord           *discord.App       `inject:"discord"`
	SettingKey        gorm.SettingKey    `inject:"settingKey"`
	SettingRepository setting.Repository `inject:"settingRepository"`
}

func (s *PresensiService) Startup() error { return nil }

func (s *PresensiService) Shutdown() error { return nil }

func (s *PresensiService) SendReminderPresensi(ctx context.Context) error {
	settings, err := s.SettingRepository.GetAllByKey(ctx, s.SettingKey.ReminderPresensi())
	if err != nil {
		return err
	}

	if len(settings) == 0 {
		return setting.ErrDataNotFound
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

		logger.Info("Send reminder presensi")
		s.Discord.Bot.ChannelMessageSend(channelID, fmt.Sprintf("Jangan lupa melakukan presensi <@&%s> !", roleID))
	}

	return nil
}
