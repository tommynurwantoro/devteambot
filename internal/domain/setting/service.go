package setting

import "context"

type Service interface {
	SetSuperAdmin(ctx context.Context, guildID, userIDs string) error
	IsSuperAdmin(ctx context.Context, guildID string, roles []string) bool
	GetPointLogChannel(ctx context.Context, guildID string) (string, error)
}
