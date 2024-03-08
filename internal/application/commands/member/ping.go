package member

import (
	"devteambot/internal/application/service"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	MessageService service.MessageService `inject:"messageService"`
}

func (c *PingCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "ping",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Check bot server status",
	}

	c.Command.AppendCommand(c.AppCommand)
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *PingCommand) Shutdown() error { return nil }

func (c *PingCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(i.Interaction)
	}
}

func (c *PingCommand) Do(i *discordgo.Interaction) {
	c.MessageService.SendStandardResponse(i, "pong", true, false)
}
