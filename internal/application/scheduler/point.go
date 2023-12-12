package scheduler

import (
	"context"
	"devteambot/internal/domain/point"
)

type PoinScheduler struct {
	Scheduler    *Scheduler    `inject:"scheduler"`
	PointService point.Service `inject:"pointService"`
}

func (s *PoinScheduler) Startup() error {
	// Every Monday 00:00
	s.Scheduler.Cron("0 0 * * 1").Do(func() {
		s.PointService.ResetQuota(context.Background())
	})

	return nil
}

func (s *PoinScheduler) Shutdown() error { return nil }
