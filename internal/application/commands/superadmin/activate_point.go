package superadmin

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ActivatePointCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	MessageService service.MessageService `inject:"messageService"`
	SettingService service.SettingService `inject:"settingService"`
}

func (c *ActivatePointCommand) Startup() error {
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *ActivatePointCommand) Shutdown() error { return nil }

func (c *ActivatePointCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "activate_point_feature") {
		c.activateButton(s, i)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "choose_thanks_channel") {
		c.chooseThanksChannel(s, i)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "choose_point_log_channel") {
		c.choosePointLogChannel(i)
	}
}

func (c *ActivatePointCommand) activateButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string

	channels, err := s.GuildChannels(i.GuildID)
	if err != nil {
		response = "Failed to get channels"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
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
		Content: "Please choose which channel you will use use thanks feature",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID: "choose_thanks_channel",
						Options:  channelOptions,
					},
				},
			},
		},
	}

	c.MessageService.SendEmbedResponse(i.Interaction, embedMessage, false)
}

func (c *ActivatePointCommand) chooseThanksChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string

	channelID := i.MessageComponentData().Values[0]

	// send embed message
	embedMessage := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Thanks",
				Description: "Use button below to send thanks",
				Color:       0x0099ff,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Send Thanks",
						Style:    discordgo.SuccessButton,
						CustomID: "send_thanks",
						Disabled: false,
					},
				},
			},
		},
	}

	if err := c.MessageService.SendEmbedMessage(channelID, embedMessage); err != nil {
		logger.Error("Failed to send embed message", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to send embed message", true, false)
		return
	}

	channels, err := s.GuildChannels(i.GuildID)
	if err != nil {
		response = "Failed to get channels"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
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
	embedMessage = &discordgo.MessageSend{
		Content: "Please choose which channel you will use to send point log",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID: "choose_point_log_channel",
						Options:  channelOptions,
					},
				},
			},
		},
	}

	c.MessageService.SendEmbedResponse(i.Interaction, embedMessage, false)

	if err := c.MessageService.DeleteMessage(i.Message.ChannelID, i.Message.ID); err != nil {
		logger.Error("Failed to delete previous message", err)
		return
	}
}

func (c *ActivatePointCommand) choosePointLogChannel(i *discordgo.InteractionCreate) {
	var response string

	channelID := i.MessageComponentData().Values[0]

	if err := c.SettingService.SetPointLogChannel(context.Background(), i.GuildID, channelID); err != nil {
		response = "Failed to set point log channel"
		logger.Error(response, err)
		c.MessageService.EditStandardResponse(i.Interaction, response)
		return
	}

	response = "Success to set point log channel"

	// delete previous message
	if err := c.MessageService.DeleteMessage(i.Message.ChannelID, i.Message.ID); err != nil {
		logger.Error("Failed to delete previous message", err)
		return
	}

	c.MessageService.SendStandardResponse(i.Interaction, response, true, false)
}
