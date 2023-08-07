package constant

type SettingKeyType uint

const (
	Admin SettingKeyType = iota
	SuperAdmin
	ReminderSholatChannel
)

type SettingKey struct {
	Key map[SettingKeyType]string
}

func NewSettingKey() SettingKey {
	key := make(map[SettingKeyType]string)
	key[Admin] = "admin"
	key[SuperAdmin] = "super_admin"
	key[ReminderSholatChannel] = "reminder_sholat_channel"

	return SettingKey{key}
}

func (c *SettingKey) Shutdown() error { return nil }

func (c *SettingKey) Admin() string {
	return c.Key[Admin]
}

func (c *SettingKey) SuperAdmin() string {
	return c.Key[SuperAdmin]
}

func (c *SettingKey) ReminderSholatChannel() string {
	return c.Key[ReminderSholatChannel]
}
