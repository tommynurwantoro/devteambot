package commands

import (
	"context"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) ThanksLeaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	ctx := context.Background()

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	categories := []string{"run", "unity", "bravery", "integrity", "customer-oriented"}

	message := &discordgo.MessageSend{
		Content: "RUBIC TOP 10 LEADERBOARD",
	}

	for _, core := range categories {
		logger.Info(core)
		embed := &discordgo.MessageEmbed{
			Title: strings.ToUpper(core),
		}

		topTen, err := c.PointRepository.GetTopTen(ctx, i.GuildID, core)
		if err != nil {
			response = "Something went wrong, can not add rubic"
			c.SendStandardResponse(i.Interaction, response, true, false)
			return
		}

		if len(topTen) == 0 {
			embed.Description = "Belum ada data"
			continue
		}

		var rank, user, rubic string
		for n, t := range topTen {
			if n == 0 {
				rank = fmt.Sprintf("%s:first_place: \n", rank)
			} else if n == 1 {
				rank = fmt.Sprintf("%s:second_place: \n", rank)
			} else if n == 2 {
				rank = fmt.Sprintf("%s:third_place: \n", rank)
			} else {
				rank = fmt.Sprintf("%s#%d\n", rank, n+1)
			}
			user = fmt.Sprintf("%s<@%s>\n", user, t.UserID)
			rubic = fmt.Sprintf("%s%d\n", rubic, t.Balance)
		}

		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Rank",
				Value:  rank,
				Inline: true,
			},
			{
				Name:   "User",
				Value:  user,
				Inline: true,
			},
			{
				Name:   "Rubic",
				Value:  rubic,
				Inline: true,
			},
		}
		embed.Color = c.Color.Green

		message.Embeds = append(message.Embeds, embed)
	}

	if len(message.Embeds) == 0 {
		message.Content = "Belum ada data"
	}

	_, err := s.ChannelMessageSendComplex(i.ChannelID, message)
	if err != nil {
		logger.Error(err.Error(), err)
		response = "Failed to send embed"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to generate leaderboard"
	c.SendStandardResponse(i.Interaction, response, true, false)
}
