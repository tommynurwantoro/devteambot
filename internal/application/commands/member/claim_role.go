package member

import (
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ClaimRoleCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	MessageService service.MessageService `inject:"messageService"`
}

func (c *ClaimRoleCommand) Startup() error {
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *ClaimRoleCommand) Shutdown() error { return nil }

func (c *ClaimRoleCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "claim_role") {
		c.Do(s, i.Interaction)
	}
}

func (c *ClaimRoleCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
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
		c.MessageService.SendStandardResponse(i, response, true, false)
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
			response = "Something went wrong, please try again later"
			logger.Error(response, err)
			c.MessageService.SendStandardResponse(i, response, true, false)
			return
		}
		response = "Success to remove role"
	} else {
		if err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID); err != nil {
			response = "Something went wrong, please try again later"
			logger.Error(response, err)
			c.MessageService.SendStandardResponse(i, response, true, false)
			return
		}
		response = "Success to add role"
	}

	c.MessageService.SendStandardResponse(i, response, true, false)
}
