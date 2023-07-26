package sql

import (
	"context"
	"devteambot/internal/domain/review"
	"devteambot/internal/domain/sharedkernel/identity"
	"devteambot/internal/pkg/logger"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB `inject:"database"`
}

func (r *ReviewRepository) Startup() error {
	logger.Info("Migrating review repository if necessary")
	err := r.DB.AutoMigrate(review.Review{})
	if err != nil {
		logger.Fatal("Error migrate review", err)
	}

	return nil
}

func (r *ReviewRepository) Shutdown() error {
	return nil
}

func (r *ReviewRepository) Create(ctx context.Context, guildID, reporter, title, url string, reviewer []string) (*review.Review, error) {
	m := review.NewReview(guildID, reporter, title, url, reviewer)

	tx := r.DB.Create(&m)
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
		return nil, tx.Error
	}

	return m, nil
}

func (r *ReviewRepository) Update(ctx context.Context, data *review.Review) error {
	tx := r.DB.Save(&data)
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
		return tx.Error
	}

	return nil
}

func (r *ReviewRepository) GetByID(ctx context.Context, id identity.ID) (*review.Review, error) {
	tx := r.DB.WithContext(ctx)

	var s = &review.Review{}
	result := tx.First(s, "id = ?", id.String())
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error get by ID: %s", result.Error.Error()), result.Error)
		return nil, result.Error
	}

	return s, nil
}

func (r *ReviewRepository) GetAllPendingByGuildID(ctx context.Context, guildID string) (review.Reviews, error) {
	reviews := make(review.Reviews, 0)

	tx := r.DB.Where("guild_id = ? AND total_pending > 0", guildID).Find(&reviews).Order("created_at")
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return reviews, nil
		}
		logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
		return nil, tx.Error
	}

	return reviews, nil
}
