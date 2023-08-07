package setting

import "context"

type Repository interface {
	GetByKey(ctx context.Context, guildID, key string, value interface{}) error
	GetAllByKey(ctx context.Context, key string) (Settings, error)
	SetValue(ctx context.Context, guildID, key string, value interface{}) error
}
