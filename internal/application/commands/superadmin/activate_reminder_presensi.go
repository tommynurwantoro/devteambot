package superadmin

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ActivateReminderPresensiCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	CacheService   cache.Service          `inject:"cache"`
	MessageService service.MessageService `inject:"messageService"`
	SettingService service.SettingService `inject:"settingService"`
}

func (c *ActivateReminderPresensiCommand) Startup() error {
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *ActivateReminderPresensiCommand) Shutdown() error { return nil }

func (c *ActivateReminderPresensiCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "activate_reminder_presensi_feature") {
		c.activateButton(s, i.Interaction)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "choose_reminder_presensi_channel") {
		c.chooseChannel(s, i.Interaction)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "choose_reminder_presensi_role") {
		c.chooseRole(i.Interaction)
	}
}

func (c *ActivateReminderPresensiCommand) activateButton(s *discordgo.Session, i *discordgo.Interaction) {
	var response string

	channels, err := s.GuildChannels(i.GuildID)
	if err != nil {
		response = "Failed to get channels"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	channelOptions := []discordgo.SelectMenuOption{}
	for _, channel := range channels {
		channelOptions = append(channelOptions, discordgo.SelectMenuOption{
			Label: channel.Name,
			Value: channel.ID,
		})
	}

	// choose channel
	embedMessage := &discordgo.MessageSend{
		Content: "Please choose which channel you will use to send reminders",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						MenuType: discordgo.ChannelSelectMenu,
						CustomID: "choose_reminder_presensi_channel",
						Options:  channelOptions,
					},
				},
			},
		},
	}

	c.MessageService.SendEmbedResponse(i, embedMessage, false)
}

func (c *ActivateReminderPresensiCommand) chooseChannel(s *discordgo.Session, i *discordgo.Interaction) {
	var response string

	channelID := i.MessageComponentData().Values[0]

	key := fmt.Sprintf("reminder_presensi_channel_%s", i.GuildID)
	if err := c.CacheService.Put(context.Background(), key, channelID, 0); err != nil {
		response = "Failed to set reminder presensi channel"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	roles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		response = "Failed to get roles"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	roleOptions := []discordgo.SelectMenuOption{}
	for _, role := range roles {
		roleOptions = append(roleOptions, discordgo.SelectMenuOption{
			Label: role.Name,
			Value: role.ID,
		})
	}

	// choose role
	content := "Please choose which role you will use to send reminder"
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType: discordgo.RoleSelectMenu,
					CustomID: "choose_reminder_presensi_role",
					Options:  roleOptions,
				},
			},
		},
	}

	c.MessageService.SendEmbedResponse(i, &discordgo.MessageSend{
		Content:    content,
		Components: components,
	}, false)

	if err := c.MessageService.DeleteMessage(i.Message.ChannelID, i.Message.ID); err != nil {
		logger.Error("Failed to delete previous message", err)
		return
	}
}

func (c *ActivateReminderPresensiCommand) chooseRole(i *discordgo.Interaction) {
	var response string
	var channelID string

	roleID := i.MessageComponentData().Values[0]

	key := fmt.Sprintf("reminder_presensi_channel_%s", i.GuildID)
	if err := c.CacheService.Get(context.Background(), key, &channelID); err != nil {
		response = "Failed to get reminder presensi channel"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	if err := c.SettingService.SetReminderPresensiChannel(context.Background(), i.GuildID, channelID, roleID); err != nil {
		response = "Failed to set reminder presensi channel"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	if _, err := c.CacheService.Delete(context.Background(), key); err != nil {
		logger.Error(response, err)
	}

	response = "Success to set reminder presensi channel"
	c.MessageService.SendStandardResponse(i, response, true, false)

	if err := c.MessageService.DeleteMessage(i.Message.ChannelID, i.Message.ID); err != nil {
		logger.Error("Failed to delete previous message", err)
		return
	}
}
