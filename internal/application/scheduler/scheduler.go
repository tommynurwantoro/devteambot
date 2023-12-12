package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	*gocron.Scheduler
}

func (s *Scheduler) Startup() error {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := gocron.NewScheduler(loc)

	s.Scheduler = scheduler

	return nil
}

func (s *Scheduler) Shutdown() error { s.Scheduler.Stop(); return nil }
