package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) Thanks(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	ctx := context.Background()

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	var to, core, reason string

	if opt, ok := optionMap["to"]; ok {
		to = opt.UserValue(s).ID
	}

	if opt, ok := optionMap["core"]; ok {
		core = opt.StringValue()
	}

	if opt, ok := optionMap["reason"]; ok {
		reason = opt.StringValue()
	}

	var pointLogChannel string
	err := c.SettingRepository.GetByKey(ctx, i.GuildID, c.SettingKey.PointLogChannel(), &pointLogChannel)
	if err != nil {
		response = "Something went wrong, please try again later"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	if pointLogChannel == "" {
		response = "Aktifkan fitur ini terlebih dahulu dengan command /activate_point_feature"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	_, err = c.PointRepository.Increase(ctx, i.GuildID, to, core, reason, 1)
	if err != nil {
		response = "Something went wrong, can not add point"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	c.SendStandardResponse(i.Interaction, "Success", true, false)

	c.App.Bot.ChannelMessageSend(pointLogChannel, fmt.Sprintf("[%s] - <@%s> barusan say thanks to <@%s> karena %s", strings.ToUpper(core), i.Member.User.ID, to, reason))
}
