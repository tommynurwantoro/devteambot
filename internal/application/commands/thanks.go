package commands

import (
	"context"
	"devteambot/internal/domain/point"
	"devteambot/internal/pkg/logger"
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

	if to == i.Member.User.ID {
		response = "Tidak bisa berterima kasih ke diri sendiri"
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	pointLogChannel, err := c.SettingService.GetPointLogChannel(ctx, i.GuildID)
	if err != nil {
		response = "Something went wrong, please try again later"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	if pointLogChannel == "" {
		response = "Aktifkan fitur ini terlebih dahulu dengan command /activate_point_feature"
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	if err := c.PointService.SendThanks(ctx, i.GuildID, i.Member.User.ID, to, core, reason); err != nil {
		if err == point.ErrLimitReached {
			response = "Limit mingguan kamu sudah habis, kamu bisa pakai command /thanks lagi mulai senin depan"
		} else {
			response = "Something went wrong, can not add rubic"
			logger.Error(response, err)
		}
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	c.MessageService.SendStandardResponse(i.Interaction, "Success", true, false)

	c.MessageService.SendStandardMessage(pointLogChannel, fmt.Sprintf("[%s] - <@%s> barusan kasih 10 rubic ke <@%s> karena %s", strings.ToUpper(core), i.Member.User.ID, to, reason))
}
