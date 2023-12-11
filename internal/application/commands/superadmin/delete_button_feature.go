package superadmin

import (
	"context"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) DeleteButtonFeature(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	var messageID string
	var index int64

	if opt, ok := optionMap["message_id"]; ok {
		messageID = opt.StringValue()
	}

	if opt, ok := optionMap["index"]; ok {
		index = opt.IntValue()
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

		for i, component := range actionRow.Components {
			if int64(i+1) == index {
				continue
			}

			components = append(components, component)
		}
	}

	if len(components) > 0 {
		message.Components = append(message.Components, discordgo.ActionsRow{
			Components: components,
		})
	}

	_, err = s.ChannelMessageEditComplex(message)
	if err != nil {
		logger.Error(err.Error(), err)
		response = "Failed to delete button"
		c.Command.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to delete button"
	c.Command.SendStandardResponse(i.Interaction, response, true, false)
}
