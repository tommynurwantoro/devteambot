package commands

import (
	"context"
	"fmt"

	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) SendStandardResponse(i *discordgo.Interaction, response string, isPrivate, isRemovePreview bool) {
	data := &discordgo.InteractionResponseData{
		Content: response,
	}

	if isPrivate {
		data.Flags = discordgo.MessageFlagsEphemeral
	}

	if isRemovePreview {
		data.Flags = data.Flags | discordgo.MessageFlagsSuppressEmbeds
	}

	err := c.App.Bot.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
	if err != nil {
		logger.Error(err.Error(), err)
	}
}

func (c *Command) EditResponse(i *discordgo.Interaction, response string) {
	_, err := c.App.Bot.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &response,
	})
	if err != nil {
		logger.Error(err.Error(), err)
	}
}

func (c *Command) Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (c *Command) AssignRole(s *discordgo.Session, i *discordgo.InteractionCreate, announceChannelID, roleID string, isAnnounce bool) {
	var response string

	err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID)
	if err != nil {
		response = "Something went wrong, please try again later"
		logger.Error(fmt.Sprintf("%s: %s", response, err.Error()), err)
		c.EditResponse(i.Interaction, response)
		return
	}

	response = fmt.Sprintf("Role <@&%s> assigned successfully", roleID)
	if isAnnounce {
		s.ChannelMessageSendComplex(announceChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("<@%s>", i.Member.User.ID),
			Embeds: []*discordgo.MessageEmbed{
				{
					Color:       c.Color.Green,
					Description: response,
					Author: &discordgo.MessageEmbedAuthor{
						Name:    i.Member.User.Username,
						IconURL: i.Member.AvatarURL(""),
					},
				},
			},
		})
	}

	c.EditResponse(i.Interaction, response)
}

func (c *Command) IsAdmin(ctx context.Context, i *discordgo.Interaction) bool {
	isAdmin := false
	admins := make([]string, 0)

	err := c.SettingRepository.GetByKey(ctx, i.GuildID, "admin", &admins)
	if err != nil {
		logger.Error("Error get setting "+err.Error(), err)
		return false
	}

	roles := i.Member.Roles
	for _, role := range roles {
		if c.Contains(admins, role) {
			isAdmin = true
			break
		}
	}

	return isAdmin
}

func (c *Command) IsSuperAdmin(ctx context.Context, i *discordgo.Interaction) bool {
	isAdmin := false
	admins := make([]string, 0)

	err := c.SettingRepository.GetByKey(ctx, i.GuildID, "super_admin", &admins)
	if err != nil {
		logger.Error("Error get setting "+err.Error(), err)
		return false
	}

	roles := i.Member.Roles
	for _, role := range roles {
		if c.Contains(admins, role) {
			isAdmin = true
			break
		}
	}

	return isAdmin
}
