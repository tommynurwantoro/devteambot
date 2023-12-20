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
