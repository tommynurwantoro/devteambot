package marketplace

import (
	"context"
	"devteambot/internal/domain/sharedkernel/identity"
)

type Repository interface {
	GetAll(ctx context.Context, guildID string) (Marketplaces, error)
	GetAllAvailable(ctx context.Context, guildID string) (Marketplaces, error)
	GetByID(ctx context.Context, id identity.ID) (*Marketplace, error)
	Create(ctx context.Context, marketplace *Marketplace) error
	Update(ctx context.Context, marketplace *Marketplace) error
	Delete(ctx context.Context, guildID, item string) error
}
