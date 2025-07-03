package member

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type AntrianReviewCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	ReviewService  service.ReviewService  `inject:"reviewService"`
	MessageService service.MessageService `inject:"messageService"`
}

func (c *AntrianReviewCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "antrian_review",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Show all pending review",
	}

	c.Command.AppendCommand(c.AppCommand)
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *AntrianReviewCommand) Shutdown() error { return nil }

func (c *AntrianReviewCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(i.Interaction)
	}
}

func (c *AntrianReviewCommand) Do(i *discordgo.Interaction) {
	var response string
	ctx := context.Background()

	reviews, err := c.ReviewService.GetAntrian(ctx, i.GuildID)
	if err != nil {
		response = "Error to get antrian review"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	content, embed := c.ReviewService.PrettyAntrian(reviews)
	message := &discordgo.MessageSend{
		Content: content,
		Embeds:  []*discordgo.MessageEmbed{embed},
	}

	c.MessageService.SendEmbedResponse(i, message, false)
}
