package scheduler

import (
	"context"
	"devteambot/internal/domain/presensi"
	"time"
)

type PresensiScheduler struct {
	Scheduler       *Scheduler       `inject:"scheduler"`
	PresensiService presensi.Service `inject:"presensiService"`
}

func (s *PresensiScheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	s.Scheduler.Every(1).Day().At("07:55").Do(func() {
		now := time.Now().In(loc)
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			return
		}

		s.PresensiService.SendReminderPresensi(context.Background())
	})

	s.Scheduler.Every(1).Day().At("17:05").Do(func() {
		now := time.Now().In(loc)
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			return
		}

		s.PresensiService.SendReminderPresensi(context.Background())
	})

	return nil
}

func (s *PresensiScheduler) Shutdown() error { return nil }
