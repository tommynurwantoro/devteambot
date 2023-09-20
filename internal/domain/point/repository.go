package point

import (
	"context"
)

type Repository interface {
	Increase(ctx context.Context, guildID, userID, category, reason string, total int64) (*Point, error)
	Decrease(ctx context.Context, guildID, userID, category, reason string, total int64) (*Point, error)
	GetByUserID(ctx context.Context, guildID, userID, category string) (*Point, error)
	GetTopTen(ctx context.Context, guildID, category string) (Points, error)
}
