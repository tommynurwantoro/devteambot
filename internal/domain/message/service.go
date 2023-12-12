package message

import "context"

type Service interface {
	SendStandardMessage(ctx context.Context, message string) error
}
