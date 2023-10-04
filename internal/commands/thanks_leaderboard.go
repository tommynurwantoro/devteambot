package commands

import (
	"context"
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

	var core string

	if opt, ok := optionMap["core"]; ok {
		core = opt.StringValue()
	}

	topTen, err := c.PointRepository.GetTopTen(ctx, i.GuildID, core)
	if err != nil {
		response = "Something went wrong, can not add rubic"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	for n, t := range topTen {
		response = fmt.Sprintf("%s%d. <@%s> : %d rubic\n", response, n+1, t.UserID, t.Balance)
	}

	if response != "" {
		response = fmt.Sprintf("Berikut leaderboard sementara untuk Core Value %s:\n%s", strings.ToUpper(core), response)
	}

	if response == "" {
		response = "Belum ada data"
	}

	c.SendStandardResponse(i.Interaction, response, false, false)
}
