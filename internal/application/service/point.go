package service

import (
	"context"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/domain/point"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
)

type PointService struct {
	Cache           cache.Service    `inject:"cache"`
	RedisKey        redis.RedisKey   `inject:"redisKey"`
	PointRepository point.Repository `inject:"pointRepository"`
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

func (s *PointService) GetTopTen(ctx context.Context, guildID, category string) (point.Points, error) {
	topTen, err := s.PointRepository.GetTopTen(ctx, guildID, category)
	if err != nil {
		return nil, err
	}

	if len(topTen) == 0 {
		return nil, point.ErrDataNotFound
	}

	return topTen, nil
}

func (s *PointService) SendThanks(ctx context.Context, guildID, from, to, core, reason string) error {
	// Cek limit
	limit := 0
	if err := s.Cache.Get(ctx, s.RedisKey.LimitThanks(guildID, from), &limit); err != nil && err != cache.ErrNil {
		return err
	}

	if limit >= 30 {
		return point.ErrLimitReached
	}

	_, err := s.PointRepository.Increase(ctx, guildID, to, core, reason, 10)
	if err != nil {
		return err
	}

	s.Cache.Increment(ctx, s.RedisKey.LimitThanks(guildID, from), 10)

	return nil
}
