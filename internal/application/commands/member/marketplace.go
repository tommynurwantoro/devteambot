package member

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/domain/marketplace"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type MarketplaceCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	MarketplaceService service.MarketplaceService `inject:"marketplaceService"`
	MessageService     service.MessageService     `inject:"messageService"`
	SettingService     service.SettingService     `inject:"settingService"`
}

func (c *MarketplaceCommand) Startup() error {
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *MarketplaceCommand) Shutdown() error { return nil }

func (c *MarketplaceCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "marketplace_buy_item") {
		c.buyItem(i.Interaction)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "marketplace_choose_item") {
		c.chooseItem(i.Interaction)
	}
}

func (c *MarketplaceCommand) buyItem(i *discordgo.Interaction) {
	allItems, err := c.MarketplaceService.GetAllItems(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get all items", err)
		c.MessageService.SendStandardResponse(i, "Failed to get all items, please try again", true, false)
		return
	}

	items := []discordgo.SelectMenuOption{}
	for _, item := range allItems {
		items = append(items, discordgo.SelectMenuOption{
			Label: item.Item,
			Value: item.ID.String(),
		})
	}

	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    "marketplace_choose_item",
		Placeholder: "Select item",
		Options:     items,
	}

	c.MessageService.SendEmbedResponse(i, &discordgo.MessageSend{
		Content: "Please choose which item you want to buy",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					selectMenu,
				},
			},
		},
	}, true)
}

func (c *MarketplaceCommand) chooseItem(i *discordgo.Interaction) {
	selectedItem := i.MessageComponentData().Values[0]

	item, err := c.MarketplaceService.GetItem(context.Background(), selectedItem)
	if err != nil {
		logger.Error("Failed to get item", err)
		c.MessageService.SendStandardResponse(i, "Failed to get item, please try again", true, false)
		return
	}

	currentPoint, err := c.MarketplaceService.BuyItem(context.Background(), i.GuildID, i.Member.User.ID, selectedItem)
	if err != nil {
		if err == marketplace.ErrInsufficientBalance {
			c.MessageService.SendStandardResponse(i, "Kamu tidak memiliki cukup rubic", true, false)
			return
		}
		if err == marketplace.ErrOutOfStock {
			c.MessageService.SendStandardResponse(i, "Item sudah habis", true, false)
			return
		}

		logger.Error("Failed to buy item", err)
		c.MessageService.SendStandardResponse(i, "Failed to buy item, please try again", true, false)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Pembelian item berhasil",
		Description: fmt.Sprintf("Pembelian *%s* berhasil dengan harga *%d* rubic\nKamu memiliki sisa *%d* rubic", item.Item, item.Price, currentPoint.Balance),
		Color:       0x0099ff,
	}

	c.updateMarketplaceMessage(i)

	c.MessageService.SendEmbedResponse(i, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}, true)

	pointLogChannel, err := c.SettingService.GetPointLogChannel(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get point log channel", err)
	} else {
		c.MessageService.SendStandardMessage(pointLogChannel, fmt.Sprintf("<@%s> barusan membeli item *%s* dengan harga *%d* rubic", i.Member.User.ID, item.Item, item.Price))
	}
}

func (c *MarketplaceCommand) updateMarketplaceMessage(i *discordgo.Interaction) {
	// get all items
	items, err := c.MarketplaceService.GetAllItems(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get all items", err)
		c.MessageService.SendStandardResponse(i, "Failed to get all items, please try again", true, false)
		return
	}

	// updatem marketplace embed message
	channelID, messageID, err := c.SettingService.GetMarketplaceMessage(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get marketplace message", err)
		c.MessageService.SendStandardResponse(i, "Failed to get marketplace message, please try again", true, false)
		return
	}

	description := ""
	for _, item := range items {
		description += fmt.Sprintf("\n🛒 %s - `%d rubic` - *total %d*", item.Item, item.Price, item.Stock)
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
		c.MessageService.SendStandardResponse(i, "Failed to update marketplace message, please try again", true, false)
		return
	}
}
