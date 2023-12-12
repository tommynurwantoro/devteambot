package scheduler

import (
	"context"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
)

type PoinScheduler struct {
	Scheduler *Scheduler     `inject:"scheduler"`
	Cache     cache.Service  `inject:"cache"`
	RedisKey  redis.RedisKey `inject:"redisKey"`
}

func (s *PoinScheduler) Startup() error {
	// Every Monday 00:00
	s.Scheduler.Cron("0 0 * * 1").Do(func() {
		s.ResetQuota(context.Background())
	})

	return nil
}

func (s *PoinScheduler) Shutdown() error { return nil }

func (s *PoinScheduler) ResetQuota(ctx context.Context) {
	allLimit, err := s.Cache.Keys(ctx, s.RedisKey.AllLimitThanks())
	if err != nil {
		logger.Error("Error: "+err.Error(), err)
		return
	}

	for _, l := range allLimit {
		_, err := s.Cache.Delete(ctx, l)
		if err != nil {
			logger.Error("Error: "+err.Error(), err)
			return
		}
	}

	logger.Info("Limit has been reset")
}
