package sholat

import "context"

type Service interface {
	GetSholatSchedule(ctx context.Context) error
}
