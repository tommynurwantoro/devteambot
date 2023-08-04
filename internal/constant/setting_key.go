package constant

type SettingKeyType uint

const (
	Admin SettingKeyType = iota
	SuperAdmin
)

type SettingKey struct {
	Key map[SettingKeyType]string
}

func NewSettingKey() SettingKey {
	key := make(map[SettingKeyType]string)
	key[Admin] = "admin"
	key[SuperAdmin] = "super_admin"

	return SettingKey{key}
}

func (c *SettingKey) Shutdown() error { return nil }

func (c *SettingKey) Admin() string {
	return c.Key[Admin]
}

func (c *SettingKey) SuperAdmin() string {
	return c.Key[SuperAdmin]
}
