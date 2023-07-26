package review

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, guildID, reporter, title, url string, reviewer []string) (*Review, error)
	Update(ctx context.Context, data *Review) error
	GetAllPendingByGuildID(ctx context.Context, guildID string) (Reviews, error)
}
