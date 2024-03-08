package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) Ask(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	ctx := context.Background()

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	var prompt string

	if opt, ok := optionMap["prompt"]; ok {
		prompt = opt.StringValue()
	}

	c.MessageService.SendStandardResponse(i.Interaction, "I'm thinking...", false, false)

	result, err := c.AIService.GetResponse(ctx, prompt)
	if err != nil {
		response = "Something went wrong, please try again later"
		c.MessageService.EditStandardResponse(i.Interaction, response)
		return
	}

	response = fmt.Sprintf("> Prompt: %s\n", prompt)
	c.MessageService.EditStandardResponse(i.Interaction, response)

	for len(result) > 2000 {
		send := result[:2000]
		result = result[2000:]

		c.MessageService.SendStandardMessage(i.ChannelID, send)
	}

	c.MessageService.SendStandardMessage(i.ChannelID, result)

}
