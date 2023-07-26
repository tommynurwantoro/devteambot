package constant

type SettingKeyType uint

const (
	Admin SettingKeyType = iota
	SuperAdmin
	CoinEmoji
	LogCoinChannelID
	ModReportChannelID
)

type SettingKey struct {
	Key map[SettingKeyType]string
}

func NewSettingKey() SettingKey {
	key := make(map[SettingKeyType]string)
	key[Admin] = "admin"
	key[SuperAdmin] = "super_admin"
	key[CoinEmoji] = "coin_emoji"
	key[LogCoinChannelID] = "log_coin_channel_id"
	key[ModReportChannelID] = "mod_report_channel_id"

	return SettingKey{key}
}

func (c *SettingKey) Shutdown() error { return nil }

func (c *SettingKey) Admin() string {
	return c.Key[Admin]
}

func (c *SettingKey) SuperAdmin() string {
	return c.Key[SuperAdmin]
}

func (c *SettingKey) CoinEmoji() string {
	return c.Key[CoinEmoji]
}

func (c *SettingKey) LogCoinChannelID() string {
	return c.Key[LogCoinChannelID]
}

func (c *SettingKey) ModReportChannelID() string {
	return c.Key[ModReportChannelID]
}
