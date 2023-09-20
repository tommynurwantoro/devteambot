package point

import (
	"devteambot/internal/domain/sharedkernel/entity"
)

type Point struct {
	entity.Entity
	GuildID  string `gorm:"not null"`
	UserID   string `gorm:"not null"`
	Category string `gorm:"not null"`
	Balance  int64  `gorm:"not null"`
}

func NewPoint(guildID, userID, category string, balance int64) *Point {
	return &Point{
		Entity:   entity.NewEntity(),
		GuildID:  guildID,
		UserID:   userID,
		Category: category,
		Balance:  balance,
	}
}

type Points []*Point
