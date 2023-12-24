package service

import (
	"context"
	"devteambot/internal/domain/setting"
	"strings"
)

type SettingService struct {
	SettingRepository setting.Repository `inject:"settingRepository"`
}

func (s *SettingService) Startup() error { return nil }

func (s *SettingService) Shutdown() error { return nil }

func (s *SettingService) SetSuperAdmin(ctx context.Context, guildID, userIDs string) error {
	listUserID := strings.Split(userIDs, ",")

	err := s.SettingRepository.SetValue(ctx, guildID, setting.SUPER_ADMIN, listUserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SettingService) IsSuperAdmin(ctx context.Context, guildID string, roles []string) bool {
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
