package sholat

import "context"

type Service interface {
	GetSholatSchedule(ctx context.Context) error
	SendReminderSholat(ctx context.Context) error
}
