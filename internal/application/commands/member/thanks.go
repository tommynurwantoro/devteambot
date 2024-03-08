package member

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/domain/point"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ThanksCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	PointService   service.PointService   `inject:"pointService"`
	SettingService service.SettingService `inject:"settingService"`
	MessageService service.MessageService `inject:"messageService"`
}

func (c *ThanksCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "thanks",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Kamu bisa kasih rubic dengan say thanks ke anggota tim yang lain",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "to",
				Description: "Orang yang dituju",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "core",
				Description: "Pilih core value yang berhubungan",
				Type:        discordgo.ApplicationCommandOptionString,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Run",
						Value: "run",
					},
					{
						Name:  "Unity",
						Value: "unity",
					},
					{
						Name:  "Bravery",
						Value: "bravery",
					},
					{
						Name:  "Integrity",
						Value: "integrity",
					},
					{
						Name:  "Customer Oriented",
						Value: "customer-oriented",
					},
				},
				Required: true,
			},
			{
				Name:        "reason",
				Description: "Alasan kamu memberikan rubic",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	}

	c.Command.AppendCommand(c.AppCommand)
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *ThanksCommand) Shutdown() error { return nil }

func (c *ThanksCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *ThanksCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
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
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	pointLogChannel, err := c.SettingService.GetPointLogChannel(ctx, i.GuildID)
	if err != nil {
		response = "Something went wrong, please try again later"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	if pointLogChannel == "" {
		response = "Aktifkan fitur ini terlebih dahulu dengan command /activate_point_feature"
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	if err := c.PointService.SendThanks(ctx, i.GuildID, i.Member.User.ID, to, core, reason); err != nil {
		if err == point.ErrLimitReached {
			response = "Limit mingguan kamu sudah habis, kamu bisa pakai command /thanks lagi mulai senin depan"
		} else {
			response = "Something went wrong, can not add rubic"
			logger.Error(response, err)
		}
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	c.MessageService.SendStandardResponse(i, "Success", true, false)

	c.MessageService.SendStandardMessage(pointLogChannel, fmt.Sprintf("[%s] - <@%s> barusan kasih 10 rubic ke <@%s> karena %s", strings.ToUpper(core), i.Member.User.ID, to, reason))
}
