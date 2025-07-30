package service

import (
	"context"
	"devteambot/internal/domain/marketplace"
	"devteambot/internal/domain/point"
	"devteambot/internal/domain/sharedkernel/identity"
)

type MarketplaceService interface {
	AddNewItem(ctx context.Context, guildID, itemName string, price, stock int64) error
	GetAllItems(ctx context.Context, guildID string) (marketplace.Marketplaces, error)
	GetItem(ctx context.Context, ID string) (*marketplace.Marketplace, error)
	UpdateItem(ctx context.Context, itemID string, price, stock int64) error
	BuyItem(ctx context.Context, guildID, userID, itemID string) (*point.Point, error)
}

type Marketplace struct {
	Repository      marketplace.Repository `inject:"marketplaceRepository"`
	PointRepository point.Repository       `inject:"pointRepository"`
}

func (s *Marketplace) AddNewItem(ctx context.Context, guildID, itemName string, price, stock int64) error {
	return s.Repository.Create(ctx, &marketplace.Marketplace{
		GuildID: guildID,
		Item:    itemName,
		Price:   price,
		Stock:   stock,
	})
}

func (s *Marketplace) GetAllItems(ctx context.Context, guildID string) (marketplace.Marketplaces, error) {
	items, err := s.Repository.GetAll(ctx, guildID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Marketplace) GetItem(ctx context.Context, ID string) (*marketplace.Marketplace, error) {
	item, err := s.Repository.GetByID(ctx, identity.FromStringOrNil(ID))
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Marketplace) UpdateItem(ctx context.Context, id string, price, stock int64) error {
	item, err := s.Repository.GetByID(ctx, identity.FromStringOrNil(id))
	if err != nil {
		return err
	}

	item.Price = price
	item.Stock = stock

	return s.Repository.Update(ctx, item)
}

func (s *Marketplace) BuyItem(ctx context.Context, guildID, userID, itemID string) (*point.Point, error) {
	item, err := s.Repository.GetByID(ctx, identity.FromStringOrNil(itemID))
	if err != nil {
		return nil, err
	}

	if item.Stock == 0 {
		return nil, marketplace.ErrOutOfStock
	}

	currentPoint, err := s.PointRepository.GetByUserID(ctx, guildID, userID)
	if err != nil {
		return nil, err
	}

	if currentPoint.Balance < item.Price {
		return nil, marketplace.ErrInsufficientBalance
	}

	point, err := s.PointRepository.Decrease(ctx, guildID, userID, "buy item "+item.Item, item.Price)
	if err != nil {
		return nil, err
	}

	item.Stock = item.Stock - 1
	if err := s.Repository.Update(ctx, item); err != nil {
		return nil, err
	}

	return point, nil
}
