package superadmin

import (
	"context"
	"devteambot/internal/pkg/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) EditEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	ctx := context.Background()

	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name != "edit_embed" {
			return
		}

		// Admin only
		if !c.Command.IsSuperAdmin(ctx, i.Interaction) {
			response := "This command is only for super admin"
			c.Command.SendStandardResponse(i.Interaction, response, true, false)
			return
		}

		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, option := range options {
			optionMap[option.Name] = option
		}

		var messageID, title, content, thumbnail string

		if opt, ok := optionMap["message_id"]; ok {
			messageID = opt.StringValue()
		}

		if opt, ok := optionMap["title"]; ok {
			title = opt.StringValue()
		}

		if opt, ok := optionMap["content"]; ok {
			content = opt.StringValue()
			content = strings.ReplaceAll(content, "|", "\n")
		}

		if opt, ok := optionMap["thumbnail"]; ok {
			thumbnail = opt.StringValue()
		}

		existingMessage, err := s.ChannelMessage(i.ChannelID, messageID)
		if err != nil {
			logger.Error(err.Error(), err)
			response = "Something went wrong, please try again later"
			c.Command.SendStandardResponse(i.Interaction, response, true, false)
			return
		}

		embedMessage := existingMessage.Embeds[0]

		if title != "" {
			embedMessage.Title = title
		}

		if content != "" {
			embedMessage.Description = content
		}

		if thumbnail != "" {
			embedMessage.Thumbnail.URL = thumbnail
		}

		existingMessage.Embeds = []*discordgo.MessageEmbed{
			embedMessage,
		}

		editMessage := &discordgo.MessageEdit{
			Content:    &existingMessage.Content,
			Components: existingMessage.Components,
			Embeds:     existingMessage.Embeds,
			ID:         existingMessage.ID,
			Channel:    existingMessage.ChannelID,
		}

		_, err = s.ChannelMessageEditComplex(editMessage)
		if err != nil {
			logger.Error(err.Error(), err)
			response = "Failed to edit embed"
			c.Command.SendStandardResponse(i.Interaction, response, true, false)
			return
		}

		response = "Success to edit embed"
		c.Command.SendStandardResponse(i.Interaction, response, true, false)
	}
}
