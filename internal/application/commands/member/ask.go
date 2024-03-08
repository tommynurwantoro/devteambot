package member

import (
	"context"
	"devteambot/internal/application/service"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type AskCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	AIService      service.AIService      `inject:"aiService"`
	MessageService service.MessageService `inject:"messageService"`
}

func (c *AskCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "ask",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Ask the bot a question",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "prompt",
				Description: "Your question",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	}

	c.Command.AppendCommand(c.AppCommand)
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *AskCommand) Shutdown() error { return nil }

func (c *AskCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(i.Interaction)
	}
}

func (c *AskCommand) Do(i *discordgo.Interaction) {
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

	c.MessageService.SendStandardResponse(i, "I'm thinking...", false, false)

	result, err := c.AIService.GetResponse(ctx, prompt)
	if err != nil {
		response = "Something went wrong, please try again later"
		c.MessageService.EditStandardResponse(i, response)
		return
	}

	response = fmt.Sprintf("> Prompt: %s\n", prompt)
	c.MessageService.EditStandardResponse(i, response)

	for len(result) > 2000 {
		send := result[:2000]
		result = result[2000:]

		c.MessageService.SendStandardMessage(i.ChannelID, send)
	}

	c.MessageService.SendStandardMessage(i.ChannelID, result)
}
