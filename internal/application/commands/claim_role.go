package commands

import (
	"devteambot/internal/pkg/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) ClaimRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string

	customID := i.MessageComponentData().CustomID
	custom := strings.Split(customID, "|")
	if len(custom) < 2 {
		return
	}

	roleID := custom[1]

	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		response = "Something went wrong, please try again later"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	found := false
	for _, role := range member.Roles {
		if role == roleID {
			found = true
		}
	}

	if found {
		if err = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID); err != nil {
		}
		response = "Success to remove role"
	} else {
		s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID)
		response = "Success to add role"
	}

	c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
}
