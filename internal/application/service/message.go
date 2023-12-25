package service

import (
	"devteambot/internal/adapter/discord"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type MessageService struct {
	App *discord.App `inject:"discord"`
}

func (s *MessageService) SendStandardResponse(i *discordgo.Interaction, response string, isPrivate, isRemovePreview bool) {
	data := &discordgo.InteractionResponseData{
		Content: response,
	}

	if isPrivate {
		data.Flags = discordgo.MessageFlagsEphemeral
	}

	if isRemovePreview {
		data.Flags = data.Flags | discordgo.MessageFlagsSuppressEmbeds
	}

	if err := s.App.Bot.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	}); err != nil {
		logger.Error("Error to send message", err)
	}
}

func (s *MessageService) EditStandardResponse(i *discordgo.Interaction, response string) {
	if _, err := s.App.Bot.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &response,
	}); err != nil {
		logger.Error("Error to edit message", err)
	}
}

func (s *MessageService) SendEmbedResponse(i *discordgo.Interaction, content string, embed *discordgo.MessageEmbed, isPrivate bool) {
	data := &discordgo.InteractionResponseData{
		Content: content,
		Embeds:  []*discordgo.MessageEmbed{embed},
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeUsers,
			},
		},
	}

	if isPrivate {
		data.Flags = discordgo.MessageFlagsEphemeral
	}

	if err := s.App.Bot.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	}); err != nil {
		logger.Error("Error to send embed message", err)
	}
}

func (s *MessageService) SendStandardMessage(channelID, message string) {
	if _, err := s.App.Bot.ChannelMessageSend(channelID, message); err != nil {
		logger.Error("Error to send message", err)
	}
}
