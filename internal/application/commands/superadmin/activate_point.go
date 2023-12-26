package superadmin

import (
	"context"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) ActivatePoint(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	var response string

	// Admin only
	if !c.Command.SettingService.IsSuperAdmin(ctx, i.GuildID, i.Member.Roles) {
		response := "This command is only for super admin"
		c.Command.MessageService.SendStandardResponse(i.Interaction, response, true, false)
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

	if err := c.SettingService.SetPointLogChannel(ctx, i.GuildID, channelID); err != nil {
		response = "Failed to activate point feature"
		logger.Error(response, err)
		c.Command.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to activate point feature"
	c.Command.MessageService.SendStandardResponse(i.Interaction, response, true, false)
}
