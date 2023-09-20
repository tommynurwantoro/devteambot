package constant

type SettingKeyType uint

const (
	Admin SettingKeyType = iota
	SuperAdmin
	ReminderSholatChannel
	ReminderPresensiChannel
	PointLogChannel
)

type SettingKey struct {
	Key map[SettingKeyType]string
}

func NewSettingKey() SettingKey {
	key := make(map[SettingKeyType]string)
	key[Admin] = "admin"
	key[SuperAdmin] = "super_admin"
	key[ReminderSholatChannel] = "reminder_sholat_channel"
	key[ReminderPresensiChannel] = "reminder_presensi_channel"
	key[PointLogChannel] = "point_log_channel"

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

func (c *SettingKey) ReminderPresensiChannel() string {
	return c.Key[ReminderPresensiChannel]
}

func (c *SettingKey) PointLogChannel() string {
	return c.Key[PointLogChannel]
}
