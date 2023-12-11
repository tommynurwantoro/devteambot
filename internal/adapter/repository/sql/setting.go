package sql

import (
	"context"
	"devteambot/internal/adapter/repository"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const SettingKey = "setting|%s|%s"

type SettingRepository struct {
	DB *repository.Gorm `inject:"database"`
}

func (r *SettingRepository) Startup() error {
	logger.Info("[table migration] Migrating setting repository if necessary")
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
	tx := r.DB.WithContext(ctx)
	var s setting.Setting

	result := tx.First(&s, "guild_id = ? AND key = ?", guildID, key)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return setting.ErrDataNotFound
		}
		logger.Error(fmt.Sprintf("Error: %s", result.Error.Error()), result.Error)
		return result.Error
	}

	err := json.Unmarshal([]byte(s.Value), &value)
	if err != nil {
		return err
	}

	return nil
}

func (r *SettingRepository) GetAllByKey(ctx context.Context, key string) (setting.Settings, error) {
	tx := r.DB.WithContext(ctx)
	s := make(setting.Settings, 0)

	result := tx.Where("key = ?", key).Find(&s)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, setting.ErrDataNotFound
		}
		logger.Error(fmt.Sprintf("Error: %s", result.Error.Error()), result.Error)
		return nil, result.Error
	}

	return s, nil
}

func (r *SettingRepository) SetValue(ctx context.Context, guildID, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		logger.Error(fmt.Sprintf("Error: %s", err.Error()), err)
		return err
	}

	data := &setting.Setting{}
	tx := r.DB.First(data, "guild_id = ? AND key = ?", guildID, key)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			data = setting.NewSetting(guildID, key, string(bytes))
		} else {
			logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
			return tx.Error
		}
	}

	data.Value = string(bytes)
	tx = r.DB.Save(&data)
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
		return tx.Error
	}

	return nil
}
