package scheduler

import (
	"context"
	"devteambot/config"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type SholatScheduler struct {
	Scheduler     *Scheduler            `inject:"scheduler"`
	Config        *config.Config        `inject:"config"`
	SholatService service.SholatService `inject:"sholatService"`
}

func (s *SholatScheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	conf, ok := s.Config.Schedulers["sholat-get-today-schedule"]
	if ok && conf.Enable {
		s.Scheduler.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(
					gocron.NewAtTime(conf.Time.Hour, conf.Time.Minute, conf.Time.Second),
				),
			),
			gocron.NewTask(
				func() {
					s.SholatService.GetTodaySchedule(context.Background())
				},
			),
		)
		logger.Info("Sholat: Get Today Schedule is enabled")
	}

	conf, ok = s.Config.Schedulers["sholat-send-reminder"]
	if ok && conf.Enable {
		s.Scheduler.NewJob(
			gocron.DurationJob(
				1*time.Minute,
			),
			gocron.NewTask(
				func() {
					if time.Now().In(loc).Weekday() == time.Saturday || time.Now().In(loc).Weekday() == time.Sunday {
						return
					}
					s.SholatService.SendReminder(context.Background())
				},
			),
		)
		logger.Info("Sholat: Send Reminder is enabled")
	}

	return nil
}

func (s *SholatScheduler) Shutdown() error { return nil }
