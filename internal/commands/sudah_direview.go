package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) SudahDireview(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	ctx := context.Background()

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	var number int64

	if opt, ok := optionMap["number"]; ok {
		number = opt.IntValue()
	}

	all, err := c.ReviewRepository.GetAllPendingByGuildID(ctx, i.GuildID)
	if err != nil {
		response = "Error to get all pending review"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	for j := 1; j <= len(all); j++ {
		if number == int64(j) {
			updating := all[j-1]
			newReviewer := make([]string, 0)
			found := false
			for _, s := range updating.Reviewer {
				if s == i.Member.User.ID {
					found = true
					response = fmt.Sprintf("FYI buat <@%s> barusan <@%s> udah selesai review [%s](%s)", updating.Reporter, s, updating.Title, updating.Url)
					continue
				}

				newReviewer = append(newReviewer, s)
			}

			updating.Reviewer = newReviewer
			updating.TotalPending = len(newReviewer)

			err := c.ReviewRepository.Update(ctx, updating)
			if err != nil {
				response = "Error to update review"
				c.SendStandardResponse(i.Interaction, response, true, false)
				return
			}

			if !found {
				response = "Kamu bukan reviewer di antrian itu"
			}

			c.SendStandardResponse(i.Interaction, response, false, false)
			return
		}
	}

	response = fmt.Sprintf("Urutan %d gak ada!", number)
	c.SendStandardResponse(i.Interaction, response, true, false)
}
