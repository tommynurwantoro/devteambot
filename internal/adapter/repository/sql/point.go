package sql

import (
	"context"
	"devteambot/internal/domain/point"
	"devteambot/internal/domain/sharedkernel/entity"
	"devteambot/internal/pkg/logger"
	"fmt"

	"gorm.io/gorm"
)

type PointRepository struct {
	DB *gorm.DB `inject:"database"`
}

func (r *PointRepository) Startup() error {
	logger.Info("Migrating point repository if necessary")
	err := r.DB.AutoMigrate(point.Point{})
	if err != nil {
		logger.Fatal("Error migrate point", err)
	}

	return nil
}

func (r *PointRepository) Shutdown() error { return nil }

func (r *PointRepository) Increase(ctx context.Context, guildID, userID string, total int64) (*point.Point, error) {
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

	return &c, nil
}

func (r *PointRepository) Decrease(ctx context.Context, guildID, userID string, total int64) (*point.Point, error) {
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

	return &c, nil
}

func (r *PointRepository) GetByUserID(ctx context.Context, guildID string, userID string) (*point.Point, error) {
	var c point.Point

	tx := r.DB.First(&c, "guild_id = ? AND user_id = ?", guildID, userID)
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
		return nil, tx.Error
	}

	return &c, nil
}

func (r *PointRepository) GetTopTen(ctx context.Context, guildID string) (point.Points, error) {
	listPoint := make(point.Points, 0)

	tx := r.DB.Where("guild_id = ?", guildID).Order("balance DESC").Find(&listPoint).Limit(10)
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("Error: %s", tx.Error.Error()), tx.Error)
		return nil, tx.Error
	}

	return listPoint, nil
}
