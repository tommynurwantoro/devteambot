package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) AntrianReview(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	ctx := context.Background()

	reviews, err := c.ReviewRepository.GetAllPendingByGuildID(ctx, i.GuildID)
	if err != nil {
		response = "Error to get review"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	for i, r := range reviews {
		reviewer := ""
		for _, user := range r.Reviewer {
			reviewer = fmt.Sprintf("%s, <@%s>", reviewer, user)
		}

		reviewer = reviewer[2:]

		response = fmt.Sprintf("%s\n%d. [%s](%s) %s", response, i+1, r.Title, r.Url, reviewer)
	}

	if response == "" {
		response = "Gak ada antrian review nih 👍🏻"
	}

	c.SendStandardResponse(i.Interaction, response, false, true)
}
