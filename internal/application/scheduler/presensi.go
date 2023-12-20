package scheduler

import (
	"context"
	"devteambot/config"
	"devteambot/internal/domain/presensi"
	"devteambot/internal/pkg/logger"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type PresensiScheduler struct {
	Scheduler       *Scheduler       `inject:"scheduler"`
	Config          *config.Config   `inject:"config"`
	PresensiService presensi.Service `inject:"presensiService"`
}

func (s *PresensiScheduler) Startup() error {
	conf, ok := s.Config.Schedulers["presensi-send-reminder-pagi"]
	if ok && conf.Enable {
		s.Scheduler.NewJob(
			gocron.WeeklyJob(
				1,
				gocron.NewWeekdays(
					time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
				),
				gocron.NewAtTimes(
					gocron.NewAtTime(
						conf.Time.Hour,
						conf.Time.Minute,
						conf.Time.Second,
					),
				),
			),
			gocron.NewTask(s.PresensiService.SendReminder, context.Background()),
		)
		logger.Info("Presensi: Send Reminder Pagi is enabled")
	}

	conf, ok = s.Config.Schedulers["presensi-send-reminder-sore"]
	if ok && conf.Enable {
		s.Scheduler.NewJob(
			gocron.WeeklyJob(
				1,
				gocron.NewWeekdays(
					time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
				),
				gocron.NewAtTimes(
					gocron.NewAtTime(
						conf.Time.Hour,
						conf.Time.Minute,
						conf.Time.Second,
					),
				),
			),
			gocron.NewTask(s.PresensiService.SendReminder, context.Background()),
		)
		logger.Info("Presensi: Send Reminder Sore is enabled")
	}

	return nil
}

func (s *PresensiScheduler) Shutdown() error { return nil }
