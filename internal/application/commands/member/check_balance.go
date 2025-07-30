package member

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CheckBalanceCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	PointService   service.PointService   `inject:"pointService"`
	MessageService service.MessageService `inject:"messageService"`
	ThanksService  service.ThanksService  `inject:"thanksService"`
}

func (c *CheckBalanceCommand) Startup() error {
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *CheckBalanceCommand) Shutdown() error { return nil }

func (c *CheckBalanceCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "check_balance") {
		c.Do(i.Interaction)
	}
}

func (c *CheckBalanceCommand) Do(i *discordgo.Interaction) {
	balance, err := c.PointService.GetPointBalance(context.Background(), i.GuildID, i.Member.User.ID)
	if err != nil {
		logger.Error("Failed to get point", err)
		c.MessageService.SendStandardResponse(i, "Failed to get point, please try again", true, false)
		return
	}

	thanksLimit, err := c.ThanksService.ThanksLimit(context.Background(), i.GuildID, i.Member.User.ID)
	if err != nil {
		logger.Error("Failed to get thanks", err)
		c.MessageService.SendStandardResponse(i, "Failed to get thanks, please try again", true, false)
		return
	}

	c.MessageService.SendEmbedResponse(i, &discordgo.MessageSend{
		Content: "Ini saldo rubic dan thanks kamu minggu ini",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Saldo Rubic",
				Description: fmt.Sprintf("Total rubic kamu: *%d*\nSisa thanks kamu minggu ini: *%d*", balance, thanksLimit),
				Color:       0x0099ff,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: i.Member.User.AvatarURL(""),
				},
			},
		},
	}, true)
}
