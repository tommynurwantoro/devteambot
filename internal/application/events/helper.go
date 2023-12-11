package events

import (
	"github.com/bwmarrin/discordgo"
)

func (c *Event) isAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	isAdmin := false

	roles := m.Member.Roles
	for _, role := range roles {
		if c.Admins[role] == true {
			isAdmin = true
			break
		}
	}

	return isAdmin
}
