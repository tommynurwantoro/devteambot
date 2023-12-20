package setting

import "context"

type Service interface {
	SetSuperAdmin(ctx context.Context, guildID, userIDs string) error
}
