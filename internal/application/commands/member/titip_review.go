package member

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type TitipReviewCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	ReviewService  service.ReviewService  `inject:"reviewService"`
	MessageService service.MessageService `inject:"messageService"`
}

func (c *TitipReviewCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "titip_review",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Titip review ke anggota squad",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "Title",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "url",
				Description: "Url",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "reviewer_1",
				Description: "Reviewer 1",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "reviewer_2",
				Description: "Reviewer 2",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    false,
			},
			{
				Name:        "reviewer_3",
				Description: "Reviewer 3",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    false,
			},
			{
				Name:        "reviewer_4",
				Description: "Reviewer 4",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    false,
			},
			{
				Name:        "reviewer_5",
				Description: "Reviewer 5",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    false,
			},
		},
	}

	c.Command.AppendCommand(c.AppCommand)
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *TitipReviewCommand) Shutdown() error { return nil }

func (c *TitipReviewCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *TitipReviewCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	var response string
	ctx := context.Background()

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	var title, url, mentions string
	var reviewers []string

	if opt, ok := optionMap["title"]; ok {
		title = opt.StringValue()
	}
	if opt, ok := optionMap["url"]; ok {
		url = opt.StringValue()
	}
	if opt, ok := optionMap["reviewer_1"]; ok {
		userID := opt.UserValue(s).ID
		reviewers = append(reviewers, userID)
		mentions = fmt.Sprintf("<@%s>", userID)
	}
	if opt, ok := optionMap["reviewer_2"]; ok {
		userID := opt.UserValue(s).ID
		reviewers = append(reviewers, userID)
		mentions = fmt.Sprintf("%s <@%s>", mentions, userID)
	}
	if opt, ok := optionMap["reviewer_3"]; ok {
		userID := opt.UserValue(s).ID
		reviewers = append(reviewers, userID)
		mentions = fmt.Sprintf("%s <@%s>", mentions, userID)
	}
	if opt, ok := optionMap["reviewer_4"]; ok {
		userID := opt.UserValue(s).ID
		reviewers = append(reviewers, userID)
		mentions = fmt.Sprintf("%s <@%s>", mentions, userID)
	}
	if opt, ok := optionMap["reviewer_5"]; ok {
		userID := opt.UserValue(s).ID
		reviewers = append(reviewers, userID)
		mentions = fmt.Sprintf("%s <@%s>", mentions, userID)
	}

	if err := c.ReviewService.AddReviewer(ctx, i.GuildID, i.Member.User.ID, title, url, reviewers); err != nil {
		response = "Error to add review"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	response = fmt.Sprintf("<@%s> barusan titip review ya... %s tolong nanti cek 🫰🏻\n[%s](%s)", i.Member.User.ID, mentions, title, url)
	c.MessageService.SendStandardResponse(i, response, false, false)
}
