package superadmin

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type ActivateReminderPresensiCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	MessageService service.MessageService `inject:"messageService"`
	SettingService service.SettingService `inject:"settingService"`
}

func (c *ActivateReminderPresensiCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "activate_reminder_presensi",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Activate reminder presensi feature",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "channel",
				Description: "Channel for the reminder",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    true,
			},
			{
				Name:        "role",
				Description: "Reminder presensi will mention this role",
				Type:        discordgo.ApplicationCommandOptionRole,
				Required:    true,
			},
		},
	}

	c.CommandSuperAdmin.AppendCommand(c.AppCommand)
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *ActivateReminderPresensiCommand) Shutdown() error { return nil }

func (c *ActivateReminderPresensiCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *ActivateReminderPresensiCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	ctx := context.Background()
	var response string

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

	if err := c.SettingService.SetReminderPresensiChannel(ctx, i.GuildID, channelID, roleID); err != nil {
		response = "Failed to activate reminder"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	response = "Success to activate reminder"
	c.MessageService.SendStandardResponse(i, response, true, false)
}
