package gorm

import (
	"context"
	"devteambot/internal/adapter/repository"
	"devteambot/internal/domain/thanks"
	"devteambot/internal/pkg/logger"
)

type ThanksRepository struct {
	DB *repository.Gorm `inject:"database"`
}

func (r *ThanksRepository) Startup() error {
	logger.Info("Migrating thanks repository if necessary")
	err := r.DB.AutoMigrate(thanks.ThanksLog{})
	if err != nil {
		logger.Fatal("Error migrate thanks", err)
	}

	return nil
}

func (r *ThanksRepository) Shutdown() error { return nil }

func (r *ThanksRepository) Create(ctx context.Context, thanksLog *thanks.ThanksLog) error {
	return r.DB.Create(thanksLog).Error
}

func (r *ThanksRepository) GetAll(ctx context.Context, guildID string) (thanks.ThanksLogs, error) {
	var thanksLogs thanks.ThanksLogs

	tx := r.DB.Where(thanks.ThanksLog{GuildID: guildID}).Order("created_at DESC").Find(&thanksLogs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return thanksLogs, nil
}
