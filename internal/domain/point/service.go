package point

import "context"

type Service interface {
	ResetQuota(ctx context.Context) error
	GetTopTen(ctx context.Context, guildID, category string) (Points, error)
	SendThanks(ctx context.Context, guildID, from, to, core, reason string) error
}
