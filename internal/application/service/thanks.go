package service

import (
	"context"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/domain/point"
	"devteambot/internal/domain/thanks"
	"devteambot/internal/pkg/cache"
)

const (
	THANKS_LIMIT      = 30
	PER_THANKS        = 10
	PER_THANKS_SENDER = 5
)

type ThanksService interface {
	Log(ctx context.Context, guildID, userID, to, coreValue, reason string) error
	SendThanks(ctx context.Context, guildID, from, to, core, reason string) error
	ThanksLimit(ctx context.Context, guildID, userID string) (int64, error)
}

type Thanks struct {
	Cache            cache.Service     `inject:"cache"`
	RedisKey         redis.RedisKey    `inject:"redisKey"`
	ThanksRepository thanks.Repository `inject:"thanksRepository"`
	PointRepository  point.Repository  `inject:"pointRepository"`
}

func (s *Thanks) Log(ctx context.Context, guildID, userID, to, coreValue, reason string) error {
	return s.ThanksRepository.Create(ctx, &thanks.ThanksLog{
		GuildID:   guildID,
		UserID:    userID,
		To:        to,
		CoreValue: coreValue,
		Reason:    reason,
	})
}

func (s *Thanks) SendThanks(ctx context.Context, guildID, from, to, core, reason string) error {
	// Cek limit
	limit := 0
	if err := s.Cache.Get(ctx, s.RedisKey.LimitThanks(guildID, from), &limit); err != nil && err != cache.ErrNil {
		return err
	}

	if limit >= THANKS_LIMIT {
		return point.ErrLimitReached
	}

	_, err := s.PointRepository.Increase(ctx, guildID, to, reason, core, PER_THANKS)
	if err != nil {
		return err
	}

	_, err = s.PointRepository.Increase(ctx, guildID, from, "send thanks", core, PER_THANKS_SENDER)
	if err != nil {
		return err
	}

	s.Cache.Increment(ctx, s.RedisKey.LimitThanks(guildID, from), PER_THANKS)

	return nil
}

func (s *Thanks) ThanksLimit(ctx context.Context, guildID, userID string) (int64, error) {
	// Cek limit
	rubicUsed := 0
	if err := s.Cache.Get(ctx, s.RedisKey.LimitThanks(guildID, userID), &rubicUsed); err != nil && err != cache.ErrNil {
		return 0, err
	}

	return int64((THANKS_LIMIT - rubicUsed) / PER_THANKS), nil
}
