package superadmin

import (
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type AddButtonFeatureCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	MessageService service.MessageService `inject:"messageService"`
	SettingService service.SettingService `inject:"settingService"`
}

func (c *AddButtonFeatureCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "add_button_feature",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Add button feature to existing message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "command",
				Description: "Choose feature command",
				Type:        discordgo.ApplicationCommandOptionString,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Claim Role",
						Value: "claim_role",
					},
				},
				Required: true,
			},
			{
				Name:        "message_id",
				Description: "Message ID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "button1",
				Description: "RoleID|Name|IconName|IconID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "button2",
				Description: "RoleID|Name|IconName|IconID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "button3",
				Description: "RoleID|Name|IconName|IconID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "button4",
				Description: "RoleID|Name|IconName|IconID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "button5",
				Description: "RoleID|Name|IconName|IconID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}

	c.CommandSuperAdmin.AppendCommand(c.AppCommand)
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *AddButtonFeatureCommand) Shutdown() error { return nil }

func (c *AddButtonFeatureCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *AddButtonFeatureCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	var response string

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	var feature, messageID string
	var buttons []string

	if opt, ok := optionMap["command"]; ok {
		feature = opt.StringValue()
	}

	if opt, ok := optionMap["message_id"]; ok {
		messageID = opt.StringValue()
	}

	if opt, ok := optionMap["button1"]; ok {
		buttons = append(buttons, opt.StringValue())
	}

	if opt, ok := optionMap["button2"]; ok {
		buttons = append(buttons, opt.StringValue())
	}

	if opt, ok := optionMap["button3"]; ok {
		buttons = append(buttons, opt.StringValue())
	}

	if opt, ok := optionMap["button4"]; ok {
		buttons = append(buttons, opt.StringValue())
	}

	if opt, ok := optionMap["button5"]; ok {
		buttons = append(buttons, opt.StringValue())
	}

	m, err := c.CommandSuperAdmin.Discord.Bot.ChannelMessage(i.ChannelID, messageID)
	if err != nil {
		response = "Something went wrong, please try again later"
		logger.Error(response, err)
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
		if err := actionRow.UnmarshalJSON(data); err != nil {
			logger.Error(err.Error(), err)
			response = "Failed to add button"
			c.MessageService.SendStandardResponse(i, response, true, false)
			return
		}

		components = append(components, actionRow.Components...)
	}

	for _, button := range buttons {
		if len(components) == 5 {
			*message.Components = append(*message.Components, discordgo.ActionsRow{
				Components: components,
			})

			components = []discordgo.MessageComponent{}
		}

		split := strings.Split(button, "|")
		if len(split) < 4 {
			response = "Please check input format"
			c.MessageService.SendStandardResponse(i, response, true, false)
			return
		}
		id := split[0]
		name := split[1]
		iconName := split[2]
		iconID := split[3]

		components = append(components, discordgo.Button{
			Emoji: &discordgo.ComponentEmoji{
				Name:     iconName,
				ID:       iconID,
				Animated: false,
			},
			Label:    name,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("%s|%s", feature, id),
		})
	}

	if len(components) > 0 {
		*message.Components = append(*message.Components, discordgo.ActionsRow{
			Components: components,
		})
	}

	_, err = s.ChannelMessageEditComplex(message)
	if err != nil {
		logger.Error(err.Error(), err)
		response = "Failed to add button"
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	response = "Success to add button"
	c.MessageService.SendStandardResponse(i, response, true, false)
}
