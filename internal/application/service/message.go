package service

import (
	"context"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/repository/gorm"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/cache"
)

type MessageService struct {
	Cache             cache.Service      `inject:"cache"`
	SettingRepository setting.Repository `inject:"settingRepository"`
	App               *discord.App       `inject:"discord"`
	RedisKey          redis.RedisKey     `inject:"redisKey"`
	SettingKey        gorm.SettingKey    `inject:"settingKey"`
}

func (s *MessageService) SendStandardMessage(ctx context.Context, message string) error {
	return nil
}
