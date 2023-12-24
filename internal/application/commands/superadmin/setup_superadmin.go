package superadmin

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) SetupSuperAdmin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	userIDs := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	err := c.SettingService.SetSuperAdmin(context.Background(), i.GuildID, userIDs)
	if err != nil {
		c.Command.MessageService.SendStandardResponse(i.Interaction, "Something went wrong, please try again later", true, false)
		return
	}

	c.Command.MessageService.SendStandardResponse(i.Interaction, "Success to set super admin.", true, false)
}
