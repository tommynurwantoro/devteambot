package commands

import (
	"context"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) AntrianReview(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	ctx := context.Background()

	reviews, err := c.ReviewService.GetAntrian(ctx, i.GuildID)
	if err != nil {
		response = "Error to get antrian review"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	content, embed := c.ReviewService.PrettyAntrian(reviews)

	c.MessageService.SendEmbedResponse(i.Interaction, content, embed, false)
}
