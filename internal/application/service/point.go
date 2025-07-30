package service

import (
	"context"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/domain/point"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
)

type PointService interface {
	ResetQuota(ctx context.Context) error
	GetTopTen(ctx context.Context, guildID string) (point.Points, error)
	GetPointBalance(ctx context.Context, guildID, userID string) (int64, error)
}

type Point struct {
	Cache           cache.Service    `inject:"cache"`
	RedisKey        redis.RedisKey   `inject:"redisKey"`
	PointRepository point.Repository `inject:"pointRepository"`
}

func (s *Point) ResetQuota(ctx context.Context) error {
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

	allThanksThisWeek, err := s.Cache.Keys(ctx, s.RedisKey.AllThanksThisWeek())
	if err != nil {
		logger.Error("Error: "+err.Error(), err)
		return err
	}

	for _, t := range allThanksThisWeek {
		_, err := s.Cache.Delete(ctx, t)
		if err != nil {
			logger.Error("Error: "+err.Error(), err)
			return err
		}
	}

	logger.Info("Limit has been reset")
	return nil
}

func (s *Point) GetTopTen(ctx context.Context, guildID string) (point.Points, error) {
	topTen, err := s.PointRepository.GetTopTen(ctx, guildID)
	if err != nil {
		return nil, err
	}

	if len(topTen) == 0 {
		return nil, point.ErrDataNotFound
	}

	return topTen, nil
}

func (s *Point) GetPointBalance(ctx context.Context, guildID, userID string) (int64, error) {
	point, err := s.PointRepository.GetByUserID(ctx, guildID, userID)
	if err != nil {
		return 0, err
	}

	return point.Balance, nil
}
