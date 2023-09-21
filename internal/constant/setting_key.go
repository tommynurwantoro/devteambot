package constant

type SettingKeyType uint

const (
	Admin SettingKeyType = iota
	SuperAdmin
	ReminderSholat
	ReminderPresensi
	PointLogChannel
)

type SettingKey struct {
	Key map[SettingKeyType]string
}

func NewSettingKey() SettingKey {
	key := make(map[SettingKeyType]string)
	key[Admin] = "admin"
	key[SuperAdmin] = "super_admin"
	key[ReminderSholat] = "reminder_sholat"
	key[ReminderPresensi] = "reminder_presensi"
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

func (c *SettingKey) ReminderSholat() string {
	// ChannelID|RoleID
	return c.Key[ReminderSholat]
}

func (c *SettingKey) ReminderPresensi() string {
	// ChannelID|RoleID
	return c.Key[ReminderPresensi]
}

func (c *SettingKey) PointLogChannel() string {
	return c.Key[PointLogChannel]
}
