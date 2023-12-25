package commands

import (
	"context"
	"devteambot/internal/domain/review"
	"devteambot/internal/pkg/logger"
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

	all, err := c.ReviewService.GetAntrian(ctx, i.GuildID)
	if err != nil {
		response = "Error to get antrian review"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	for j := 1; j <= len(all); j++ {
		if number == int64(j) {
			updating := all[j-1]
			err := c.ReviewService.UpdateDone(ctx, updating, i.Member.User.ID)
			if err != nil {
				if err == review.ErrReviewerNotFound {
					response = "Kamu bukan reviewer di antrian itu"
				} else {
					response = "Error to update review"
					logger.Error(response, err)
				}

				c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
				return
			}

			response = fmt.Sprintf("FYI buat <@%s> barusan <@%s> udah selesai review [%s](%s)", updating.Reporter, i.Member.User.ID, updating.Title, updating.Url)
			c.MessageService.SendStandardResponse(i.Interaction, response, false, false)
			return
		}
	}

	response = fmt.Sprintf("Urutan %d gak ada!", number)
	c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
}
