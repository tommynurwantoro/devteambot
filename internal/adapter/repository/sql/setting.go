package sql

import (
	"context"
	"devteambot/internal/adapter/cache"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const SettingKey = "setting|%s|%s"

type SettingRepository struct {
	Cache cache.Cache `inject:"cache"`
	DB    *gorm.DB    `inject:"database"`
}

func (r *SettingRepository) Startup() error {
	logger.Info("Migrating setting repository if necessary")
	err := r.DB.AutoMigrate(setting.Setting{})
	if err != nil {
		logger.Fatal("Error migrate setting", err)
	}

	// r.DB.Create(&setting.Setting{
	// 	Entity:  entity.NewEntity(),
	// 	GuildID: "982659357953118208",
	// 	Key:     "admin",
	// 	Value:   "[\"983995219013947453\", \"986884584128020481\"]",
	// })

	// r.DB.Create(&setting.Setting{
	// 	Entity:  entity.NewEntity(),
	// 	GuildID: "982659357953118208",
	// 	Key:     "super_admin",
	// 	Value:   "[\"983995219013947453\"]",
	// })

	return nil
}

func (r *SettingRepository) Shutdown() error {
	return nil
}

func (r *SettingRepository) GetByKey(ctx context.Context, guildID, key string, value interface{}) error {
	var s setting.Setting
	redisKey := fmt.Sprintf(SettingKey, guildID, key)

	err := r.Cache.Get(ctx, redisKey, &s)
	if err != nil {
		if err != cache.ErrNil {
			return err
		}

		tx := r.DB.First(&s, "guild_id = ? AND key = ?", guildID, key)
		if tx.Error != nil {
			logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
			return tx.Error
		}

		r.Cache.Put(ctx, redisKey, s, 30*time.Second)
	}

	err = json.Unmarshal([]byte(s.Value), &value)
	if err != nil {
		return err
	}

	return nil
}

func (r *SettingRepository) GetAllByKey(ctx context.Context, key string) (setting.Settings, error) {
	s := make(setting.Settings, 0)
	redisKey := fmt.Sprintf(SettingKey, "all", key)

	err := r.Cache.Get(ctx, redisKey, &s)
	if err != nil {
		if err != cache.ErrNil {
			return nil, err
		}

		tx := r.DB.Where("key = ?", key).Find(&s)
		if tx.Error != nil {
			if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				return s, nil
			}
			logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
			return nil, tx.Error
		}

		r.Cache.Put(ctx, redisKey, s, 30*time.Second)
	}

	return s, nil
}
