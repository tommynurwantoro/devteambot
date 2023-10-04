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
	if err := c.SettingRepository.GetByKey(ctx, i.GuildID, c.SettingKey.PointLogChannel(), &pointLogChannel); err != nil {
		response = "Something went wrong, please try again later"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	if pointLogChannel == "" {
		response = "Aktifkan fitur ini terlebih dahulu dengan command /activate_point_feature"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	var limit int64
	if err := c.Cache.Get(ctx, c.RedisKey.LimitThanks(i.GuildID, i.Member.User.ID), &limit); err != nil {
		response = "Something went wrong, please try again later"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	if limit >= 30 {
		response = "Limit mingguan kamu sudah habis, kamu bisa pakai command /thanks lagi mulai senin depan"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	if _, err := c.PointRepository.Increase(ctx, i.GuildID, to, core, reason, 10); err != nil {
		response = "Something went wrong, can not add rubic"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	c.Cache.Increment(ctx, c.RedisKey.LimitThanks(i.GuildID, i.Member.User.ID), 10)

	c.SendStandardResponse(i.Interaction, "Success", true, false)

	c.App.Bot.ChannelMessageSend(pointLogChannel, fmt.Sprintf("[%s] - <@%s> barusan kasih 10 rubic ke <@%s> karena %s", strings.ToUpper(core), i.Member.User.ID, to, reason))
}
