package superadmin

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type NewItemData struct {
	Name  string
	Price int64
}

type MarketplaceCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	CacheService       cache.Service              `inject:"cache"`
	MarketplaceService service.MarketplaceService `inject:"marketplaceService"`
	MessageService     service.MessageService     `inject:"messageService"`
	SettingService     service.SettingService     `inject:"settingService"`
}

func (c *MarketplaceCommand) Startup() error {
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *MarketplaceCommand) Shutdown() error { return nil }

func (c *MarketplaceCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "marketplace_new_item") {
		c.newItem(i)
	} else if i.Type == discordgo.InteractionModalSubmit && strings.HasPrefix(i.ModalSubmitData().CustomID, "marketplace_submit_new_item") {
		c.submitNewItem(i)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "marketplace_update_item") {
		c.updateItem(i)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "choose_marketplace_item") {
		c.chooseUpdateItem(i)
	} else if i.Type == discordgo.InteractionModalSubmit && strings.HasPrefix(i.ModalSubmitData().CustomID, "marketplace_submit_update_item") {
		c.submitUpdateItem(i)
	}
}

func (c *MarketplaceCommand) newItem(i *discordgo.InteractionCreate) {
	modal := discordgo.InteractionResponseData{
		Title:    "Add New Item",
		CustomID: "marketplace_submit_new_item",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "marketplace_item_name",
						Label:       "Item Name",
						Style:       discordgo.TextInputShort,
						Placeholder: "Write your item name here",
						Required:    true,
						MaxLength:   100,
						MinLength:   1,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "marketplace_item_price",
						Label:       "Price",
						Style:       discordgo.TextInputShort,
						Placeholder: "Write your item price here",
						Required:    true,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "marketplace_item_stock",
						Label:       "Stock",
						Style:       discordgo.TextInputShort,
						Placeholder: "Write your item stock here",
						Required:    true,
					},
				},
			},
		},
	}

	c.MessageService.SendModalResponse(i.Interaction, &modal)
}

func (c *MarketplaceCommand) submitNewItem(i *discordgo.InteractionCreate) {
	itemName := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	itemPrice := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	itemStock := i.ModalSubmitData().Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	price, err := strconv.ParseInt(itemPrice, 10, 64)
	if err != nil {
		logger.Error("Failed to parse price", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to save your choice, please try again", true, false)
		return
	}

	stock, err := strconv.ParseInt(itemStock, 10, 64)
	if err != nil {
		logger.Error("Failed to parse stock", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to save your choice, please try again", true, false)
		return
	}
	if err := c.MarketplaceService.AddNewItem(context.Background(), i.GuildID, itemName, price, stock); err != nil {
		logger.Error("Failed to add new item", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to add new item, please try again", true, false)
		return
	}

	c.updateMarketplaceMessage(i)

	c.MessageService.SendStandardResponse(i.Interaction, "New item added successfully", true, false)
}

func (c *MarketplaceCommand) updateItem(i *discordgo.InteractionCreate) {
	allItems, err := c.MarketplaceService.GetAllItems(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get all items", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to get all items, please try again", true, false)
		return
	}

	items := []discordgo.SelectMenuOption{}
	for _, item := range allItems {
		items = append(items, discordgo.SelectMenuOption{
			Label: item.Item,
			Value: item.ID.String(),
		})
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType: discordgo.StringSelectMenu,
					CustomID: "choose_marketplace_item",
					Options:  items,
				},
			},
		},
	}

	c.MessageService.SendEmbedResponse(i.Interaction, &discordgo.MessageSend{
		Content:    "Please choose which item you will update",
		Components: components,
	}, true)
}

func (c *MarketplaceCommand) chooseUpdateItem(i *discordgo.InteractionCreate) {
	selectedItem := i.MessageComponentData().Values[0]

	modal := discordgo.InteractionResponseData{
		Title:    "Update Item",
		CustomID: fmt.Sprintf("marketplace_submit_update_item|%s", selectedItem),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "marketplace_item_price",
						Label:       "Price",
						Style:       discordgo.TextInputShort,
						Placeholder: "Write your item price here",
						Required:    true,
						MaxLength:   100,
						MinLength:   1,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "marketplace_item_stock",
						Label:       "Stock",
						Style:       discordgo.TextInputShort,
						Placeholder: "Write your item stock here",
						Required:    true,
					},
				},
			},
		},
	}

	c.MessageService.SendModalResponse(i.Interaction, &modal)
}

func (c *MarketplaceCommand) submitUpdateItem(i *discordgo.InteractionCreate) {
	selectedItem := strings.Split(i.ModalSubmitData().CustomID, "|")[1]
	itemPrice := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	itemStock := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	price, err := strconv.ParseInt(itemPrice, 10, 64)
	if err != nil {
		logger.Error("Failed to parse price", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to save your choice, please try again", true, false)
		return
	}

	stock, err := strconv.ParseInt(itemStock, 10, 64)
	if err != nil {
		logger.Error("Failed to parse stock", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to save your choice, please try again", true, false)
		return
	}

	if err := c.MarketplaceService.UpdateItem(context.Background(), selectedItem, price, stock); err != nil {
		logger.Error("Failed to update item", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to update item, please try again", true, false)
		return
	}

	c.updateMarketplaceMessage(i)

	c.MessageService.SendStandardResponse(i.Interaction, "Item updated successfully", true, false)
}

func (c *MarketplaceCommand) updateMarketplaceMessage(i *discordgo.InteractionCreate) {
	// get all items
	items, err := c.MarketplaceService.GetAllItems(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get all items", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to get all items, please try again", true, false)
		return
	}

	// updatem marketplace embed message
	channelID, messageID, err := c.SettingService.GetMarketplaceMessage(context.Background(), i.GuildID)
	if err != nil {
		logger.Error("Failed to get marketplace message", err)
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to get marketplace message, please try again", true, false)
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
		c.MessageService.SendStandardResponse(i.Interaction, "Failed to update marketplace message, please try again", true, false)
		return
	}
}
