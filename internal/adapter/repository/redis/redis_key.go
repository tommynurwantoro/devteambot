package redis

import "fmt"

type KeyType uint

const (
	DailySholatSchedule KeyType = iota
	LimitThanks
	AllLimitThanks
	AllThanksThisWeek
)

type RedisKey struct {
	Key map[KeyType]string
}

func NewRedisKey() RedisKey {
	key := make(map[KeyType]string)
	key[DailySholatSchedule] = "daily_sholat_schedule"
	key[LimitThanks] = "limit_thanks|%s|%s"
	key[AllLimitThanks] = "limit_thanks|*"
	key[AllThanksThisWeek] = "thanks_this_week|*"

	return RedisKey{key}
}

func (c *RedisKey) Shutdown() error { return nil }

func (c *RedisKey) DailySholatSchedule() string {
	return fmt.Sprintf(c.Key[DailySholatSchedule])
}

func (c *RedisKey) LimitThanks(guildID, userID string) string {
	return fmt.Sprintf(c.Key[LimitThanks], guildID, userID)
}

func (c *RedisKey) AllLimitThanks() string {
	return c.Key[AllLimitThanks]
}

func (c *RedisKey) AllThanksThisWeek() string {
	return c.Key[AllThanksThisWeek]
}
