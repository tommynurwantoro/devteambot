package scheduler

import (
	"context"
	"devteambot/internal/domain/sholat"
	"devteambot/internal/pkg/logger"
	"time"
)

type SholatScheduler struct {
	Scheduler     *Scheduler     `inject:"scheduler"`
	SholatService sholat.Service `inject:"sholatService"`
}

func (s *SholatScheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	s.Scheduler.Every(1).Day().At("03:00").Do(func() {
		logger.Info("Get Sholat Schedule")
		s.SholatService.GetSholatSchedule(context.Background())
	})

	s.Scheduler.Every(1).Minute().Do(func() {
		now := time.Now().In(loc)
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			return
		}
		s.SholatService.SendReminderSholat(context.Background())
	})

	return nil
}

func (s *SholatScheduler) Shutdown() error { return nil }
