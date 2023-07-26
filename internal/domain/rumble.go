package domain

type RumbleAttr struct {
	ChannelID string
	MessageID string
	Reward    int64
	Entries   []RumbleEntry
	Winner    RumbleEntry
}

type RumbleEntry struct {
	ID   string
	Name string
}
