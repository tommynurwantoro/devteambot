package presensi

import "context"

type Service interface {
	SendReminder(ctx context.Context) error
}
