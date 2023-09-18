package superadmin

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) ActivateReminderPresensi(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	var response string

	// Admin only
	if !c.Command.IsSuperAdmin(ctx, i.Interaction) {
		response := "This command is only for super admin"
		c.Command.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	var channelID string

	if opt, ok := optionMap["channel_id"]; ok {
		channelID = opt.ChannelValue(s).ID
	}

	err := c.SettingRepository.SetValue(ctx, i.GuildID, c.Command.SettingKey.ReminderPresensiChannel(), channelID)
	if err != nil {
		response = "Failed to activate reminder"
		c.Command.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to activate reminder"
	c.Command.SendStandardResponse(i.Interaction, response, true, false)
}
