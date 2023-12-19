package scheduler

import (
	"context"
	"devteambot/config"
	"devteambot/internal/domain/sholat"
	"devteambot/internal/pkg/logger"
	"time"
)

type SholatScheduler struct {
	Scheduler     *Scheduler     `inject:"scheduler"`
	Config        *config.Config `inject:"config"`
	SholatService sholat.Service `inject:"sholatService"`
}

func (s *SholatScheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	conf, ok := s.Config.Schedulers["sholat-get-today-schedule"]
	if ok && conf.Enable {
		s.Scheduler.Every(1).Day().At(conf.Time).Do(func() {
			s.SholatService.GetTodaySchedule(context.Background())
		})
		logger.Info("Sholat: Get Today Schedule is enabled")
	}

	conf, ok = s.Config.Schedulers["sholat-send-reminder"]
	if ok && conf.Enable {
		s.Scheduler.Every(1).Minute().Do(func() {
			now := time.Now().In(loc)
			if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
				return
			}
			s.SholatService.SendReminder(context.Background())
		})
		logger.Info("Sholat: Send Reminder is enabled")
	}

	return nil
}

func (s *SholatScheduler) Shutdown() error { return nil }
