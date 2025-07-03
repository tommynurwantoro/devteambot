package superadmin

import (
	"devteambot/internal/adapter/discord"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type DeleteButtonFeatureCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command     `inject:"commandSuperAdmin"`
	Discord           *discord.App `inject:"discord"`

	MessageService service.MessageService `inject:"messageService"`
	SettingService service.SettingService `inject:"settingService"`
}

func (c *DeleteButtonFeatureCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "delete_button_feature",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Delete button to existing message claim role",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "message_id",
				Description: "Message ID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "index",
				Description: "First button is 1",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    true,
			},
		},
	}

	c.CommandSuperAdmin.AppendCommand(c.AppCommand)
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *DeleteButtonFeatureCommand) Shutdown() error { return nil }

func (c *DeleteButtonFeatureCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *DeleteButtonFeatureCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	var response string

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	var messageID string
	var index int64

	if opt, ok := optionMap["message_id"]; ok {
		messageID = opt.StringValue()
	}

	if opt, ok := optionMap["index"]; ok {
		index = opt.IntValue()
	}

	m, err := c.Discord.Bot.ChannelMessage(i.ChannelID, messageID)
	if err != nil {
		response = "Something went wrong, please try again later"
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	message := &discordgo.MessageEdit{
		Embeds:     &m.Embeds,
		Components: &[]discordgo.MessageComponent{},
		ID:         messageID,
		Channel:    i.ChannelID,
	}

	components := []discordgo.MessageComponent{}

	for _, comp := range m.Components {
		if len(components) == 5 {
			*message.Components = append(*message.Components, discordgo.ActionsRow{
				Components: components,
			})

			components = []discordgo.MessageComponent{}
		}

		data, _ := comp.MarshalJSON()

		actionRow := discordgo.ActionsRow{}
		if err = actionRow.UnmarshalJSON(data); err != nil {
			logger.Error(err.Error(), err)
			response = "Failed to delete button"
			c.MessageService.SendStandardResponse(i, response, true, false)
			return
		}

		for i, component := range actionRow.Components {
			if int64(i+1) == index {
				continue
			}

			components = append(components, component)
		}
	}

	if len(components) > 0 {
		*message.Components = append(*message.Components, discordgo.ActionsRow{
			Components: components,
		})
	}

	_, err = s.ChannelMessageEditComplex(message)
	if err != nil {
		logger.Error(err.Error(), err)
		response = "Failed to delete button"
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	response = "Success to delete button"
	c.MessageService.SendStandardResponse(i, response, true, false)
}
