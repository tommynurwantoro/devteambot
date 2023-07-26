package setting

import (
	"devteambot/internal/domain/sharedkernel/entity"
)

type Setting struct {
	entity.Entity
	GuildID string `gorm:"not null"`
	Key     string `gorm:"not null"`
	Value   string `gorm:"not null"`
}
