package superadmin

import (
	"context"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) ActivateReminderSholat(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	var channelID, roleID string

	if opt, ok := optionMap["channel"]; ok {
		channelID = opt.ChannelValue(s).ID
	}

	if opt, ok := optionMap["role"]; ok {
		roleID = opt.RoleValue(s, i.GuildID).ID
	}

	if err := c.SettingService.SetReminderSholatChannel(ctx, i.GuildID, channelID, roleID); err != nil {
		response = "Failed to activate reminder"
		logger.Error(response, err)
		c.Command.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to activate reminder"
	c.Command.MessageService.SendStandardResponse(i.Interaction, response, true, false)
}
