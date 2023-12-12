package presensi

import "context"

type Service interface {
	SendReminderPresensi(ctx context.Context) error
}
