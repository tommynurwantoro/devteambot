package service

import (
	"context"
	"devteambot/internal/domain/setting"
	"fmt"
	"strings"
)

type SettingService interface {
	SetSuperAdmin(ctx context.Context, guildID, userIDs string) error
	IsSuperAdmin(ctx context.Context, guildID string, roles []string) bool
	GetPointLogChannel(ctx context.Context, guildID string) (string, error)
	SetPointLogChannel(ctx context.Context, guildID, channelID string) error
	SetReminderPresensiChannel(ctx context.Context, guildID, channelID, roleID string) error
	SetReminderSholatChannel(ctx context.Context, guildID, channelID, roleID string) error
}

type Setting struct {
	SettingRepository setting.Repository `inject:"settingRepository"`
}

func (s *Setting) Startup() error { return nil }

func (s *Setting) Shutdown() error { return nil }

func (s *Setting) SetSuperAdmin(ctx context.Context, guildID, userIDs string) error {
	listUserID := strings.Split(userIDs, ",")

	err := s.SettingRepository.SetValue(ctx, guildID, setting.SUPER_ADMIN, listUserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Setting) IsSuperAdmin(ctx context.Context, guildID string, roles []string) bool {
	isAdmin := false
	admins := make([]string, 0)

	err := s.SettingRepository.GetByKey(ctx, guildID, setting.SUPER_ADMIN, &admins)
	if err != nil {
		return false
	}

	if len(admins) == 0 {
		return false
	}

	for _, role := range roles {
		if contains(admins, role) {
			isAdmin = true
			break
		}
	}

	return isAdmin
}

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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
