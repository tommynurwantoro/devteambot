package scheduler

import (
	"context"
	"devteambot/config"
	"devteambot/internal/domain/point"
	"devteambot/internal/pkg/logger"
)

type PoinScheduler struct {
	Scheduler    *Scheduler     `inject:"scheduler"`
	Config       *config.Config `inject:"config"`
	PointService point.Service  `inject:"pointService"`
}

func (s *PoinScheduler) Startup() error {
	conf, ok := s.Config.Schedulers["point-reset-quota"]
	if ok && conf.Enable {
		// Every Monday 00:00
		s.Scheduler.Cron(conf.Time).Do(func() {
			s.PointService.ResetQuota(context.Background())
		})
		logger.Info("Point: Reset Quota is enabled")
	}

	return nil
}

func (s *PoinScheduler) Shutdown() error { return nil }
