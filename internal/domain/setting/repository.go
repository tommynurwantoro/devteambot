package setting

import "context"

type Repository interface {
	GetByKey(ctx context.Context, guildID, key string, value interface{}) error
}
