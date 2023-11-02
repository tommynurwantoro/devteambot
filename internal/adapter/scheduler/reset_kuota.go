package scheduler

import (
	"context"
	"devteambot/internal/pkg/logger"
)

func (s *Scheduler) ResetQuota(ctx context.Context) {
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
