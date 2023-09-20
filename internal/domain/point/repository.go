package point

import (
	"context"
)

type Repository interface {
	Increase(ctx context.Context, guildID, userID string, total int64) (*Point, error)
	Decrease(ctx context.Context, guildID, userID string, total int64) (*Point, error)
	GetByUserID(ctx context.Context, guildID, userID string) (*Point, error)
	GetTopTen(ctx context.Context, guildID string) (Points, error)
}
