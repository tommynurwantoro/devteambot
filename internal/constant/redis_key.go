package constant

import "fmt"

type KeyType uint

const (
	DailySholatSchedule KeyType = iota
	LimitThanks
)

type RedisKey struct {
	Key map[KeyType]string
}

func NewRedisKey() RedisKey {
	key := make(map[KeyType]string)
	key[DailySholatSchedule] = "daily-sholat-schedule"
	key[LimitThanks] = "limit-thanks-%s-%s"

	return RedisKey{key}
}

func (c *RedisKey) Shutdown() error { return nil }

func (c *RedisKey) DailySholatSchedule() string {
	return fmt.Sprintf(c.Key[DailySholatSchedule])
}

func (c *RedisKey) LimitThanks(guildID, userID string) string {
	return fmt.Sprintf(c.Key[LimitThanks], guildID, userID)
}
