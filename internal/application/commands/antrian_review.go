package commands

import (
	"context"
	"devteambot/internal/pkg/constant"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

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

	if len(reviews) == 0 {
		response = "Gak ada antrian review nih 👍🏻"
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Antrian Review",
		Description: "Reviewer dapat menggunakan command `/sudah_direview` untuk menyelesaikan review. Berikut adalah antrian review tim kamu:",
		Color:       constant.GREEN,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	content := ""

	for i, r := range reviews {
		reviewer := ""
		for _, user := range r.Reviewer {
			if !strings.Contains(content, user) {
				content = fmt.Sprintf("%s <@%s>", content, user)
			}
			reviewer = fmt.Sprintf("%s, <@%s>", reviewer, user)
		}

		reviewer = reviewer[2:]

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%d. %s :arrow_right: %s", i+1, r.Title, r.Url),
			Value: fmt.Sprintf("Reviewer: %s", reviewer),
		})
	}

	c.MessageService.SendEmbedResponse(i.Interaction, content, embed, false)
}
