package point

import (
	"devteambot/internal/domain/sharedkernel/entity"
	"devteambot/internal/domain/sharedkernel/identity"
)

type Point struct {
	entity.Entity
	GuildID string `gorm:"not null"`
	UserID  string `gorm:"not null"`
	Balance int64  `gorm:"not null"`
}

func NewPoint(guildID, userID string, balance int64) *Point {
	return &Point{
		Entity:  entity.NewEntity(),
		GuildID: guildID,
		UserID:  userID,
		Balance: balance,
	}
}

type Points []*Point

type PointHistory struct {
	entity.Entity
	PointID  identity.ID `gorm:"not null"`
	Reason   string      `gorm:"not null"`
	Changes  int64       `gorm:"not null"`
	Category string      `gorm:"not null;default:'run'"`
}

type PointHistories []*PointHistory

func Categories() []string {
	return []string{"run", "unity", "bravery", "integrity", "customer-oriented"}
}
