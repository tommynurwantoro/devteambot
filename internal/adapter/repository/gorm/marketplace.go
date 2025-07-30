package gorm

import (
	"context"
	"devteambot/internal/adapter/repository"
	"devteambot/internal/domain/marketplace"
	"devteambot/internal/domain/sharedkernel/entity"
	"devteambot/internal/domain/sharedkernel/identity"
	"devteambot/internal/pkg/logger"
)

type MarketplaceRepository struct {
	DB *repository.Gorm `inject:"database"`
}

func (r *MarketplaceRepository) Startup() error {
	logger.Info("Migrating marketplace repository if necessary")
	err := r.DB.AutoMigrate(marketplace.Marketplace{})
	if err != nil {
		logger.Fatal("Error migrate marketplace", err)
	}

	return nil
}

func (r *MarketplaceRepository) Shutdown() error { return nil }

func (r *MarketplaceRepository) GetAll(ctx context.Context, guildID string) (marketplace.Marketplaces, error) {
	var c marketplace.Marketplaces

	tx := r.DB.Where(marketplace.Marketplace{GuildID: guildID}).Order("created_at desc").Find(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return c, nil
}

func (r *MarketplaceRepository) GetAllAvailable(ctx context.Context, guildID string) (marketplace.Marketplaces, error) {
	var c marketplace.Marketplaces

	tx := r.DB.Where(marketplace.Marketplace{GuildID: guildID}).Where("stock > ?", 0).Find(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return c, nil
}

func (r *MarketplaceRepository) GetByID(ctx context.Context, id identity.ID) (*marketplace.Marketplace, error) {
	var c marketplace.Marketplace

	tx := r.DB.Where(marketplace.Marketplace{Entity: entity.Entity{ID: id}}).First(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &c, nil
}

func (r *MarketplaceRepository) Create(ctx context.Context, marketplace *marketplace.Marketplace) error {
	tx := r.DB.Create(&marketplace)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *MarketplaceRepository) Update(ctx context.Context, marketplace *marketplace.Marketplace) error {
	tx := r.DB.Save(&marketplace)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *MarketplaceRepository) Delete(ctx context.Context, guildID, item string) error {
	tx := r.DB.Where(marketplace.Marketplace{GuildID: guildID, Item: item}).Delete(&marketplace.Marketplace{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
