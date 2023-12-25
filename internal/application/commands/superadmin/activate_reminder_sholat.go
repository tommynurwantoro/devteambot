package superadmin

import (
	"context"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"
	"fmt"

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

	if err := c.SettingRepository.SetValue(ctx, i.GuildID, setting.REMINDER_SHOLAT, fmt.Sprintf("%s|%s", channelID, roleID)); err != nil {
		response = "Failed to activate reminder"
		logger.Error(response, err)
		c.Command.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to activate reminder"
	c.Command.MessageService.SendStandardResponse(i.Interaction, response, true, false)
}
