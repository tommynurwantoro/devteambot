package thanks

import "context"

type Repository interface {
	Create(ctx context.Context, thanksLog *ThanksLog) error
	GetAll(ctx context.Context, guildID string) (ThanksLogs, error)
}
