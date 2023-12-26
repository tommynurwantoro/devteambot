package setting

import "context"

type Service interface {
	SetSuperAdmin(ctx context.Context, guildID, userIDs string) error
	IsSuperAdmin(ctx context.Context, guildID string, roles []string) bool
	GetPointLogChannel(ctx context.Context, guildID string) (string, error)
	SetPointLogChannel(ctx context.Context, guildID, channelID string) error
	SetReminderPresensiChannel(ctx context.Context, guildID, channelID, roleID string) error
	SetReminderSholatChannel(ctx context.Context, guildID, channelID, roleID string) error
}
