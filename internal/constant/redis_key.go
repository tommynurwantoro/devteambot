package constant

import "fmt"

type KeyType uint

const (
	BingoStarted KeyType = iota
	GamificationBalance
	GamificationBattle
	GamificationCollect
	GiveawayActive
	GiveawayActives
	GiveawayAttr
	GiveawayTotal
	ModCount
	LimiterWLQuota
	RandomCoinCounter
	ReportMod
	RumbleAttr
	RumbleStarted
	WalletClaimed
	WalletGuaranteedClaimed
	VerificationCodeCounter
	VerificationCodeQuota
)

type RedisKey struct {
	Key map[KeyType]string
}

func NewRedisKey() RedisKey {
	key := make(map[KeyType]string)
	key[BingoStarted] = "bingo-%s-started"
	key[GamificationBalance] = "gamification-balance-%s-%s"
	key[GamificationBattle] = "gamification-battle-%s-%s"
	key[GamificationCollect] = "gamification-collect-%s-%s"
	key[GiveawayActive] = "ga-active-%s"
	key[GiveawayActives] = "ga-active"
	key[GiveawayAttr] = "ga-attr-%s"
	key[GiveawayTotal] = "ga-total"
	key[ModCount] = "mod-count-%s-%s"
	key[RandomCoinCounter] = "randomcoincounter-%s"
	key[ReportMod] = "report-mod-%s"
	key[VerificationCodeCounter] = "verification-code-counter-%s"
	key[VerificationCodeQuota] = "verification-code-quota"

	return RedisKey{key}
}

func (c *RedisKey) Shutdown() error { return nil }

func (c *RedisKey) BingoStarted(guildID string) string {
	return fmt.Sprintf(c.Key[BingoStarted], guildID)
}

func (c *RedisKey) GamificationBalance(guildID, username string) string {
	return fmt.Sprintf(c.Key[GamificationBalance], guildID, username)
}

func (c *RedisKey) GamificationBattle(guildID, username string) string {
	return fmt.Sprintf(c.Key[GamificationBattle], guildID, username)
}

func (c *RedisKey) GamificationCollect(guildID, username string) string {
	return fmt.Sprintf(c.Key[GamificationCollect], guildID, username)
}

func (c *RedisKey) GiveawayActive(id string) string {
	return fmt.Sprintf(c.Key[GiveawayActive], id)
}

func (c *RedisKey) GiveawayActives() string {
	return c.Key[GiveawayActives]
}

func (c *RedisKey) GiveawayAttr(id string) string {
	return fmt.Sprintf(c.Key[GiveawayAttr], id)
}

func (c *RedisKey) GiveawayTotal() string {
	return c.Key[GiveawayTotal]
}

func (c *RedisKey) ModCount(guildID, userID string) string {
	return fmt.Sprintf(c.Key[ModCount], guildID, userID)
}

func (c *RedisKey) RandomCoinCounter(guildID string) string {
	return fmt.Sprintf(c.Key[RandomCoinCounter], guildID)
}

func (c *RedisKey) ReportMod(date string) string {
	return fmt.Sprintf(c.Key[ReportMod], date)
}

func (c *RedisKey) VerificationCodeCounter(code string) string {
	return fmt.Sprintf(c.Key[VerificationCodeCounter], code)
}

func (c *RedisKey) VerificationCodeQuota() string {
	return fmt.Sprintf(c.Key[VerificationCodeQuota])
}
