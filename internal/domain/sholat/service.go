package sholat

import "context"

type Service interface {
	GetTodaySchedule(ctx context.Context) error
	SendReminder(ctx context.Context) error
}
