// Package entity represents a shared kernel for domain model
package entity

import (
	"devteambot/internal/domain/sharedkernel/identity"
	"time"

	"gopkg.in/guregu/null.v4"
)

// Entity represents domain Entity
type Entity struct {
	ID        identity.ID `json:"id" gorm:"primaryKey;type:uuid;default:public.gen_random_uuid()"`
	CreatedAt time.Time   `json:"createdAt" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time   `json:"updatedAt" gorm:"not null;autoUpdateTime"`
	DeletedAt null.Time   `json:"deletedAt"`
}

func NewEntity() Entity {
	now := time.Now()
	return Entity{
		ID:        identity.NewID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
