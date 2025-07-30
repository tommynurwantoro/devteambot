package marketplace

import (
	"devteambot/internal/domain/sharedkernel/entity"
)

type Marketplace struct {
	entity.Entity
	GuildID string `gorm:"not null"`
	Item    string `gorm:"not null"`
	Price   int64  `gorm:"not null"`
	Stock   int64  `gorm:"not null"`
}

func NewMarketplace(guildID, channelID, item string, price, stock int64) *Marketplace {
	return &Marketplace{
		Entity:  entity.NewEntity(),
		GuildID: guildID,
		Item:    item,
		Price:   price,
		Stock:   stock,
	}
}

type Marketplaces []*Marketplace
