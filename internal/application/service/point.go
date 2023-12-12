package service

import (
	"context"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
)

type PointService struct {
	Cache    cache.Service  `inject:"cache"`
	RedisKey redis.RedisKey `inject:"redisKey"`
}

func (s *PointService) ResetQuota(ctx context.Context) error {
	allLimit, err := s.Cache.Keys(ctx, s.RedisKey.AllLimitThanks())
	if err != nil {
		logger.Error("Error: "+err.Error(), err)
		return err
	}

	for _, l := range allLimit {
		_, err := s.Cache.Delete(ctx, l)
		if err != nil {
			logger.Error("Error: "+err.Error(), err)
			return err
		}
	}

	logger.Info("Limit has been reset")
	return nil
}
