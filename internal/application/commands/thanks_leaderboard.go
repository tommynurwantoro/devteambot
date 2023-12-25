package commands

import (
	"context"
	"devteambot/internal/domain/point"
	"devteambot/internal/pkg/constant"
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

	message := &discordgo.MessageSend{
		Content: "RUBIC TOP 10 LEADERBOARD",
	}

	for _, core := range point.Categories() {
		embed := &discordgo.MessageEmbed{
			Title: strings.ToUpper(core),
		}

		topTen, err := c.PointService.GetTopTen(ctx, i.GuildID, core)
		if err != nil && err != point.ErrDataNotFound {
			response = "Something went wrong, can not add rubic"
			logger.Error(response, err)
			c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
			return
		}

		if len(topTen) == 0 {
			embed.Description = "Belum ada data"
			continue
		}

		for n, t := range topTen {
			if n == 0 {
				embed.Description = fmt.Sprintf("%s:first_place: <@%s> `Total Rubic: %d`\n", embed.Description, t.UserID, t.Balance)
			} else if n == 1 {
				embed.Description = fmt.Sprintf("%s:second_place: <@%s> `Total Rubic: %d`\n", embed.Description, t.UserID, t.Balance)
			} else if n == 2 {
				embed.Description = fmt.Sprintf("%s:third_place: <@%s> `Total Rubic: %d`\n", embed.Description, t.UserID, t.Balance)
			} else {
				embed.Description = fmt.Sprintf("%s#%d <@%s> `Total Rubic: %d`\n", embed.Description, n+1, t.UserID, t.Balance)
			}
		}

		embed.Color = constant.GREEN
		message.Embeds = append(message.Embeds, embed)
	}

	if len(message.Embeds) == 0 {
		message.Content = "Belum ada data"
	}

	_, err := s.ChannelMessageSendComplex(i.ChannelID, message)
	if err != nil {
		response = "Failed to send embed"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to generate leaderboard"
	c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
}
