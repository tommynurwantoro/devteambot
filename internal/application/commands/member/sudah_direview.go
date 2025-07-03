package member

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/domain/review"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type SudahDireviewCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	ReviewService  service.ReviewService  `inject:"reviewService"`
	MessageService service.MessageService `inject:"messageService"`
}

func (c *SudahDireviewCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "sudah_direview",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Sudah melakukan review",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "number",
				Description: "Antrian nomor berapa",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    true,
			},
		},
	}

	c.Command.AppendCommand(c.AppCommand)
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *SudahDireviewCommand) Shutdown() error { return nil }

func (c *SudahDireviewCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(i.Interaction)
	}
}

func (c *SudahDireviewCommand) Do(i *discordgo.Interaction) {
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
		c.MessageService.SendStandardResponse(i, response, true, false)
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

				c.MessageService.SendStandardResponse(i, response, true, false)
				return
			}

			response = fmt.Sprintf("FYI buat <@%s> barusan <@%s> udah selesai review [%s](%s)", updating.Reporter, i.Member.User.ID, updating.Title, updating.Url)
			reviews, err := c.ReviewService.GetAntrian(ctx, i.GuildID)
			if err != nil {
				response = "Error to get antrian review"
				logger.Error(response, err)
			}

			_, embed := c.ReviewService.PrettyAntrian(reviews)
			message := &discordgo.MessageSend{
				Content: response,
				Embeds:  []*discordgo.MessageEmbed{embed},
			}

			c.MessageService.SendEmbedResponse(i, message, false)
			return
		}
	}

	response = fmt.Sprintf("Urutan %d gak ada!", number)
	c.MessageService.SendStandardResponse(i, response, true, false)
}
