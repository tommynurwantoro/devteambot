package point

import "context"

type Service interface {
	ResetQuota(ctx context.Context) error
}
