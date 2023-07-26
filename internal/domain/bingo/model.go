package bingo

import "devteambot/internal/domain/sharedkernel/entity"

type Bingo struct {
	entity.Entity
	ChannelID    string `gorm:"not null"`
	RandomNumber int64  `gorm:"not null"`
	Reward       int64  `gorm:"not null"`
}
