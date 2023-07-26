package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) TitipReview(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	_, err := c.ReviewRepository.Create(ctx, i.GuildID, i.Member.User.ID, title, url, reviewers)
	if err != nil {
		response = "Error to add review"
		c.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = fmt.Sprintf("<@%s> barusan titip review ya... %s tolong nanti cek 🫰🏻", i.Member.User.ID, mentions)
	c.SendStandardResponse(i.Interaction, response, false, false)
}
