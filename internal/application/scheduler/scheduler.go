package scheduler

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	gocron.Scheduler
}

func (s *Scheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
	)
	if err != nil {
		return err
	}

	s.Scheduler = scheduler

	return nil
}

func (s *Scheduler) Shutdown() error { s.Scheduler.Shutdown(); return nil }
