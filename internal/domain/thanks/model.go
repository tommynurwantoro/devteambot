package thanks

import "devteambot/internal/domain/sharedkernel/entity"

type ThanksLog struct {
	entity.Entity
	GuildID   string
	UserID    string
	To        string
	CoreValue string
	Reason    string
}

type ThanksLogs []*ThanksLog
