package superadmin

import (
	"context"
	"strings"

	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func (c *CommandSuperAdmin) SendEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	var response string

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

	var title, content, thumbnail string

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

	message := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       title,
				Color:       c.Command.Color.Blue,
				Description: content,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: thumbnail,
				},
			},
		},
	}

	_, err := s.ChannelMessageSendComplex(i.ChannelID, message)

	if err != nil {
		logger.Error(err.Error(), err)
		response = "Failed to send embed"
		c.Command.SendStandardResponse(i.Interaction, response, true, false)
		return
	}

	response = "Success to send embed"
	c.Command.SendStandardResponse(i.Interaction, response, true, false)
}
