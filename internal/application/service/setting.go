package service

import (
	"context"
	"devteambot/internal/domain/setting"
	"errors"
	"fmt"
	"strings"
)

type SettingService interface {
	GetPointLogChannel(ctx context.Context, guildID string) (string, error)
	SetPointLogChannel(ctx context.Context, guildID, channelID string) error
	SetReminderPresensiChannel(ctx context.Context, guildID, channelID, roleID string) error
	SetReminderSholatChannel(ctx context.Context, guildID, channelID, roleID string) error
	SetMarketplaceMessage(ctx context.Context, guildID, channelID, mesageID string) error
	GetMarketplaceMessage(ctx context.Context, guildID string) (string, string, error)
}

type Setting struct {
	SettingRepository setting.Repository `inject:"settingRepository"`
}

func (s *Setting) Startup() error { return nil }

func (s *Setting) Shutdown() error { return nil }

func (s *Setting) GetPointLogChannel(ctx context.Context, guildID string) (string, error) {
	var channelID string

	err := s.SettingRepository.GetByKey(ctx, guildID, setting.POINT_LOG_CHANNEL, &channelID)
	if err != nil {
		return "", err
	}

	return channelID, nil
}

func (s *Setting) SetPointLogChannel(ctx context.Context, guildID, channelID string) error {
	if err := s.SettingRepository.SetValue(ctx, guildID, setting.POINT_LOG_CHANNEL, channelID); err != nil {
		return err
	}

	return nil
}

func (s *Setting) SetReminderPresensiChannel(ctx context.Context, guildID, channelID, roleID string) error {
	value := fmt.Sprintf("%s|%s", channelID, roleID)
	if err := s.SettingRepository.SetValue(ctx, guildID, setting.REMINDER_PRESENSI, value); err != nil {
		return err
	}

	return nil
}

func (s *Setting) SetReminderSholatChannel(ctx context.Context, guildID, channelID, roleID string) error {
	value := fmt.Sprintf("%s|%s", channelID, roleID)
	if err := s.SettingRepository.SetValue(ctx, guildID, setting.REMINDER_SHOLAT, value); err != nil {
		return err
	}

	return nil
}

func (s *Setting) SetMarketplaceMessage(ctx context.Context, guildID, channelID, messageID string) error {
	value := fmt.Sprintf("%s|%s", channelID, messageID)
	if err := s.SettingRepository.SetValue(ctx, guildID, setting.MARKETPLACE_MESSAGE, value); err != nil {
		return err
	}

	return nil
}

func (s *Setting) GetMarketplaceMessage(ctx context.Context, guildID string) (string, string, error) {
	var value string

	err := s.SettingRepository.GetByKey(ctx, guildID, setting.MARKETPLACE_MESSAGE, &value)
	if err != nil {
		return "", "", err
	}

	split := strings.Split(value, "|")
	if len(split) != 2 {
		return "", "", errors.New("invalid marketplace message")
	}

	return split[0], split[1], nil
}
