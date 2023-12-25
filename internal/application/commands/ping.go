package commands

import (
	"github.com/bwmarrin/discordgo"
)

func (c *Command) Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.MessageService.SendStandardResponse(i.Interaction, "pong", true, false)
}
