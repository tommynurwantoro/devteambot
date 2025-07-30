package gorm

import (
	"context"
	"devteambot/internal/adapter/repository"
	"devteambot/internal/domain/point"
	"devteambot/internal/domain/sharedkernel/entity"
	"devteambot/internal/pkg/logger"
	"errors"

	"gorm.io/gorm"
)

type PointRepository struct {
	DB *repository.Gorm `inject:"database"`
}

func (r *PointRepository) Startup() error {
	logger.Info("Migrating point repository if necessary")
	err := r.DB.AutoMigrate(point.Point{}, point.PointHistory{})
	if err != nil {
		logger.Fatal("Error migrate point", err)
	}

	return nil
}

func (r *PointRepository) Shutdown() error { return nil }

func (r *PointRepository) Increase(ctx context.Context, guildID, userID, reason, category string, total int64) (*point.Point, error) {
	var c point.Point

	tx := r.DB.Where(point.Point{GuildID: guildID, UserID: userID}).
		Attrs(point.Point{Entity: entity.NewEntity(), Balance: total}).
		FirstOrCreate(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		c.Balance = c.Balance + total
		r.DB.Save(&c)
	}

	history := &point.PointHistory{
		Entity:   entity.NewEntity(),
		PointID:  c.ID,
		Reason:   reason,
		Changes:  total,
		Category: category,
	}
	tx = r.DB.Create(&history)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &c, nil
}

func (r *PointRepository) Decrease(ctx context.Context, guildID, userID, reason string, total int64) (*point.Point, error) {
	var c point.Point

	tx := r.DB.Where(point.Point{GuildID: guildID, UserID: userID}).
		Attrs(point.Point{Entity: entity.NewEntity(), Balance: 0}).
		FirstOrCreate(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		c.Balance = c.Balance - total
		r.DB.Save(&c)
	}

	history := &point.PointHistory{
		Entity:  entity.NewEntity(),
		PointID: c.ID,
		Reason:  reason,
		Changes: -total,
	}
	tx = r.DB.Create(&history)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &c, nil
}

func (r *PointRepository) GetByUserID(ctx context.Context, guildID, userID string) (*point.Point, error) {
	var c point.Point

	tx := r.DB.First(&c, "guild_id = ? AND user_id = ?", guildID, userID)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return &point.Point{
				Entity:  entity.NewEntity(),
				GuildID: guildID,
				UserID:  userID,
				Balance: 0,
			}, nil
		}
		return nil, tx.Error
	}

	return &c, nil
}

func (r *PointRepository) GetTopTen(ctx context.Context, guildID string) (point.Points, error) {
	listPoint := make(point.Points, 0)

	tx := r.DB.Where("guild_id = ? AND balance > 0", guildID).Order("balance DESC").Limit(10).Find(&listPoint)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return listPoint, nil
}
