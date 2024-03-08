package service

import (
	"context"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
)

type PresensiService interface {
	SendReminder(ctx context.Context) error
}

type Presensi struct {
	Discord           *discord.App       `inject:"discord"`
	SettingRepository setting.Repository `inject:"settingRepository"`
}

func (s *Presensi) Startup() error { return nil }

func (s *Presensi) Shutdown() error { return nil }

func (s *Presensi) SendReminder(ctx context.Context) error {
	settings, err := s.SettingRepository.GetAllByKey(ctx, setting.REMINDER_PRESENSI)
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
