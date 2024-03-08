package superadmin

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type ActivatePointCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	MessageService service.MessageService `inject:"messageService"`
	SettingService service.SettingService `inject:"settingService"`
}

func (c *ActivatePointCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "activate_point_feature",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Activate point feature",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "channel_id",
				Description: "Channel",
				Type:        discordgo.ApplicationCommandOptionChannel,
				Required:    true,
			},
		},
	}

	c.CommandSuperAdmin.AppendCommand(c.AppCommand)
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *ActivatePointCommand) Shutdown() error { return nil }

func (c *ActivatePointCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *ActivatePointCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	ctx := context.Background()
	var response string

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
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	response = "Success to activate point feature"
	c.MessageService.SendStandardResponse(i, response, true, false)
}
