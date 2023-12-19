package scheduler

import (
	"context"
	"devteambot/config"
	"devteambot/internal/domain/presensi"
	"devteambot/internal/pkg/logger"
	"time"
)

type PresensiScheduler struct {
	Scheduler       *Scheduler       `inject:"scheduler"`
	Config          *config.Config   `inject:"config"`
	PresensiService presensi.Service `inject:"presensiService"`
}

func (s *PresensiScheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	conf, ok := s.Config.Schedulers["presensi-send-reminder-pagi"]
	if ok && conf.Enable {
		s.Scheduler.Every(1).Day().At(conf.Time).Do(func() {
			now := time.Now().In(loc)
			if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
				return
			}

			s.PresensiService.SendReminder(context.Background())
		})
		logger.Info("Presensi: Send Reminder Pagi is enabled")
	}

	conf, ok = s.Config.Schedulers["presensi-send-reminder-sore"]
	if ok && conf.Enable {
		s.Scheduler.Every(1).Day().At(conf.Time).Do(func() {
			now := time.Now().In(loc)
			if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
				return
			}

			s.PresensiService.SendReminder(context.Background())
		})
		logger.Info("Presensi: Send Reminder Sore is enabled")
	}

	return nil
}

func (s *PresensiScheduler) Shutdown() error { return nil }
