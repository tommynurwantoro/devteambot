package superadmin

import (
	"context"
	"fmt"
	"strings"

	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) AddButtonFeature(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	m, err := c.Command.Discord.Bot.ChannelMessage(i.ChannelID, messageID)
	if err != nil {
		response = "Something went wrong, please try again later"
		c.Command.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	message := &discordgo.MessageEdit{
		Embeds:     m.Embeds,
		Components: []discordgo.MessageComponent{},
		ID:         messageID,
		Channel:    i.ChannelID,
	}

	components := []discordgo.MessageComponent{}

	for _, c := range m.Components {
		if len(components) == 5 {
			message.Components = append(message.Components, discordgo.ActionsRow{
				Components: components,
			})

			components = []discordgo.MessageComponent{}
		}

		data, _ := c.MarshalJSON()

		actionRow := discordgo.ActionsRow{}
		actionRow.UnmarshalJSON(data)

		components = append(components, actionRow.Components...)
	}

	for _, button := range buttons {
		if len(components) == 5 {
			message.Components = append(message.Components, discordgo.ActionsRow{
				Components: components,
			})

			components = []discordgo.MessageComponent{}
		}

		split := strings.Split(button, "|")
		if len(split) < 4 {
			response = "Please check input format"
			c.Command.SendStandardResponse(i.Interaction, response, true, false)
			return
		}
		id := split[0]
		name := split[1]
		iconName := split[2]
		iconID := split[3]

		components = append(components, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
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
		message.Components = append(message.Components, discordgo.ActionsRow{
			Components: components,
		})
	}

	_, err = s.ChannelMessageEditComplex(message)
	if err != nil {
		logger.Error(err.Error(), err)
		response = "Failed to add button"
		c.Command.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to add button"
	c.Command.SendStandardResponse(i.Interaction, response, true, false)
}
