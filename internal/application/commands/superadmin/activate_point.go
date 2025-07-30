package superadmin

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ActivatePointCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	MarketplaceService service.MarketplaceService `inject:"marketplaceService"`
	MessageService     service.MessageService     `inject:"messageService"`
	SettingService     service.SettingService     `inject:"settingService"`
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
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "choose_marketplace_channel") {
		c.chooseMarketplaceChannel(s, i)
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
						MenuType: discordgo.ChannelSelectMenu,
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
						Emoji: &discordgo.ComponentEmoji{
							Name: "👍",
						},
						Label:    "Send Thanks",
						Style:    discordgo.SuccessButton,
						CustomID: "send_thanks",
						Disabled: false,
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "💰",
						},
						Label:    "Check Balance",
						Style:    discordgo.SuccessButton,
						CustomID: "check_balance",
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
		Content: "Please choose which channel you will use to use marketplace feature",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						MenuType: discordgo.ChannelSelectMenu,
						CustomID: "choose_marketplace_channel",
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

func (c *ActivatePointCommand) chooseMarketplaceChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string

	channelID := i.MessageComponentData().Values[0]

	// send embed message
	embedMessage := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Marketplace",
				Description: "No Item Available",
				Color:       0x0099ff,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "🛒",
						},
						Label:    "Buy Item",
						Style:    discordgo.SuccessButton,
						CustomID: "marketplace_buy_item",
						Disabled: false,
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "💰",
						},
						Label:    "Check Balance",
						Style:    discordgo.SuccessButton,
						CustomID: "check_balance",
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

	// get latest message in that channel
	messages, err := s.ChannelMessages(channelID, 1, "", "", "")
	if err != nil {
		logger.Error("Failed to get latest message", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to get latest message", true, false)
		return
	}

	// set marketplace message to setting
	err = c.SettingService.SetMarketplaceMessage(context.Background(), i.GuildID, channelID, messages[0].ID)
	if err != nil {
		logger.Error("Failed to set marketplace message", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to set marketplace message", true, false)
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
						MenuType: discordgo.ChannelSelectMenu,
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

	c.updateMarketplaceMessage(i)
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

func (c *ActivatePointCommand) updateMarketplaceMessage(i *discordgo.InteractionCreate) {
	// get all items
	items, err := c.MarketplaceService.GetAllItems(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get all items", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to get all items, please try again", true, false)
		return
	}

	if len(items) == 0 {
		return
	}

	// update marketplace embed message
	channelID, messageID, err := c.SettingService.GetMarketplaceMessage(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get marketplace message", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to get marketplace message, please try again", true, false)
		return
	}

	description := ""
	for _, item := range items {
		description += fmt.Sprintf("\n🛒 %s - `harga: %d rubic` - `total stock: %d`", item.Item, item.Price, item.Stock)
	}

	messageEdit := &discordgo.MessageEdit{
		ID:      messageID,
		Channel: channelID,
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Marketplace",
				Description: description,
				Color:       0x0099ff,
			},
		},
	}

	if err := c.MessageService.EditEmbedMessage(messageEdit); err != nil {
		logger.Error("Failed to update marketplace message", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to update marketplace message, please try again", true, false)
		return
	}
}
