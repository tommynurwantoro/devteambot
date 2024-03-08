package service

import (
	"devteambot/internal/adapter/discord"
	"devteambot/internal/domain/message"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type MessageService interface {
	SendStandardResponse(i *discordgo.Interaction, response string, isPrivate, isRemovePreview bool)
	EditStandardResponse(i *discordgo.Interaction, response string)
	SendEmbedResponse(i *discordgo.Interaction, content string, embed *discordgo.MessageEmbed, isPrivate bool)
	SendStandardMessage(channelID, message string) error
	SendEmbedMessage(channelID string, message *discordgo.MessageSend) error
	EditEmbedMessage(message *discordgo.MessageEdit) error
	DeleteMessage(channelID, messageID string) error
	SendDMMessage(userID, m string) error
}

type Message struct {
	App *discord.App `inject:"discord"`
}

func (s *Message) SendStandardResponse(i *discordgo.Interaction, response string, isPrivate, isRemovePreview bool) {
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

func (s *Message) EditStandardResponse(i *discordgo.Interaction, response string) {
	if _, err := s.App.Bot.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &response,
	}); err != nil {
		logger.Error("Error to edit message", err)
	}
}

func (s *Message) SendEmbedResponse(i *discordgo.Interaction, content string, embed *discordgo.MessageEmbed, isPrivate bool) {
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

func (s *Message) SendStandardMessage(channelID, m string) error {
	if _, err := s.App.Bot.ChannelMessageSend(channelID, m); err != nil {
		logger.Error("Error to send message", err)
		return message.ErrFailedToSendMessage
	}

	return nil
}

func (s *Message) SendEmbedMessage(channelID string, m *discordgo.MessageSend) error {
	if _, err := s.App.Bot.ChannelMessageSendComplex(channelID, m); err != nil {
		logger.Error("Error to send embed message", err)
		return message.ErrFailedToSendMessage
	}

	return nil
}

func (s *Message) EditEmbedMessage(m *discordgo.MessageEdit) error {
	if _, err := s.App.Bot.ChannelMessageEditComplex(m); err != nil {
		logger.Error("Error to edit message", err)
		return message.ErrFailedToSendMessage
	}

	return nil
}

func (s *Message) DeleteMessage(channelID, messageID string) error {
	if err := s.App.Bot.ChannelMessageDelete(channelID, messageID); err != nil {
		logger.Error("Error to delete message", err)
		return message.ErrFailedToDeleteMessage
	}

	return nil
}

func (s *Message) SendDMMessage(userID, m string) error {
	channel, err := s.App.Bot.UserChannelCreate(userID)
	if err != nil {
		logger.Error("Error to create DM channel", err)
		return message.ErrCreatePrivateChat
	}

	if _, err := s.App.Bot.ChannelMessageSend(channel.ID, m); err != nil {
		logger.Error("Error to send DM message", err)
		return message.ErrSendPrivateChat
	}

	return nil
}
