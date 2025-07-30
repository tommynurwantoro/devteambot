package member

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/domain/point"
	"devteambot/internal/pkg/constant"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ThanksLeaderboardCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	PointService   service.PointService   `inject:"pointService"`
	MessageService service.MessageService `inject:"messageService"`
}

func (c *ThanksLeaderboardCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "thanks_leaderboard",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Show rubic leaderboard per category",
	}

	c.Command.AppendCommand(c.AppCommand)
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *ThanksLeaderboardCommand) Shutdown() error { return nil }

func (c *ThanksLeaderboardCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(i.Interaction)
	}
}

func (c *ThanksLeaderboardCommand) Do(i *discordgo.Interaction) {
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

	embed := &discordgo.MessageEmbed{
		Title: "RUBIC TOP 10 LEADERBOARD",
	}

	topTen, err := c.PointService.GetTopTen(ctx, i.GuildID)
	if err != nil && err != point.ErrDataNotFound {
		response = "Something went wrong, can not add rubic"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	if len(topTen) == 0 {
		embed.Description = "Belum ada data"
	} else {
		for n, t := range topTen {
			switch n {
			case 0:
				embed.Description = fmt.Sprintf("%s:first_place: <@%s> `Total Rubic: %d`\n", embed.Description, t.UserID, t.Balance)
			case 1:
				embed.Description = fmt.Sprintf("%s:second_place: <@%s> `Total Rubic: %d`\n", embed.Description, t.UserID, t.Balance)
			case 2:
				embed.Description = fmt.Sprintf("%s:third_place: <@%s> `Total Rubic: %d`\n", embed.Description, t.UserID, t.Balance)
			default:
				embed.Description = fmt.Sprintf("%s#%d <@%s> `Total Rubic: %d`\n", embed.Description, n+1, t.UserID, t.Balance)
			}
		}
	}

	embed.Color = constant.GREEN
	message.Embeds = append(message.Embeds, embed)

	if len(message.Embeds) == 0 {
		message.Content = "Belum ada data"
	}

	if err := c.MessageService.SendEmbedMessage(i.ChannelID, message); err != nil {
		response = "Failed to send embed"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	response = "Success to generate leaderboard"
	c.MessageService.SendStandardResponse(i, response, true, false)
}
