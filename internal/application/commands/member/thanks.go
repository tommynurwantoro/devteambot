package member

import (
	"context"
	"devteambot/internal/application/service"
	"devteambot/internal/domain/point"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ThanksData struct {
	To        string
	CoreValue string
	Reason    string
}

type ThanksCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	CacheService   cache.Service          `inject:"cache"`
	SettingService service.SettingService `inject:"settingService"`
	MessageService service.MessageService `inject:"messageService"`
	ThanksService  service.ThanksService  `inject:"thanksService"`
}

func (c *ThanksCommand) Startup() error {
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *ThanksCommand) Shutdown() error { return nil }

func (c *ThanksCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "send_thanks") {
		c.sendThanks(i.Interaction)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "thanks_choose_user") {
		c.chooseUser(i.Interaction)
	} else if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(i.MessageComponentData().CustomID, "thanks_choose_core_value") {
		c.chooseCoreValue(i.Interaction)
	} else if i.Type == discordgo.InteractionModalSubmit && strings.HasPrefix(i.ModalSubmitData().CustomID, "thanks_modal_reason") {
		c.sendThanksReason(i.Interaction)
	}
}

func (c *ThanksCommand) sendThanks(i *discordgo.Interaction) {
	thanksLimit, err := c.ThanksService.ThanksLimit(context.Background(), i.GuildID, i.Member.User.ID)
	if err != nil {
		logger.Error("Failed to get thanks limit", err)
		c.MessageService.SendStandardResponse(i, "Failed to get thanks limit, please try again", true, false)
		return
	}

	if thanksLimit <= 0 {
		c.MessageService.SendStandardResponse(i, "Kamu sudah mencapai limit thanks minggu ini, kamu bisa gunakan thanks lagi mulai senin depan", true, false)
		return
	}

	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.UserSelectMenu,
		CustomID:    "thanks_choose_user",
		Placeholder: "Select user",
		Options:     []discordgo.SelectMenuOption{},
	}

	c.MessageService.SendEmbedResponse(i, &discordgo.MessageSend{
		Content: "Please choose which user you want to send point",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					selectMenu,
				},
			},
		},
	}, true)
}

func (c *ThanksCommand) chooseUser(i *discordgo.Interaction) {
	to := i.MessageComponentData().Values[0]

	if to == i.Member.User.ID {
		c.MessageService.SendStandardResponse(i, "Tidak bisa berterima kasih ke diri sendiri", true, false)
		return
	}

	// check if user already in cache
	key := fmt.Sprintf("thanks_this_week|%s|%s|%s", i.GuildID, i.Member.User.ID, to)
	exists, err := c.CacheService.Exists(context.Background(), key)
	if err != nil {
		logger.Error("Failed to check if user already in cache", err)
		c.MessageService.SendStandardResponse(i, "Failed to check if user already in cache, please try again", true, false)
		return
	}

	if exists {
		c.MessageService.SendStandardResponse(i, "Kamu sudah berterima kasih ke user ini minggu ini", true, false)
		return
	}

	thanksData := ThanksData{
		To: to,
	}

	key = fmt.Sprintf("thanks|%s|%s", i.GuildID, i.Member.User.ID)
	if err := c.CacheService.Put(context.Background(), key, thanksData, 0); err != nil {
		logger.Error("Failed to put cache", err)
		c.MessageService.SendStandardResponse(i, "Failed to save your choice, please try again", true, false)
		return
	}

	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.StringSelectMenu,
		CustomID:    "thanks_choose_core_value",
		Placeholder: "Select core value",
		Options: []discordgo.SelectMenuOption{
			{
				Label: "Run",
				Value: "run",
			},
			{
				Label: "Unity",
				Value: "unity",
			},
			{
				Label: "Bravery",
				Value: "bravery",
			},
			{
				Label: "Integrity",
				Value: "integrity",
			},
			{
				Label: "Customer Oriented",
				Value: "customer-oriented",
			},
		},
	}

	c.MessageService.SendEmbedResponse(i, &discordgo.MessageSend{
		Content: "Please choose core value that indicate your thanks",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					selectMenu,
				},
			},
		},
	}, true)
}

func (c *ThanksCommand) chooseCoreValue(i *discordgo.Interaction) {
	var thanksData ThanksData
	coreValue := i.MessageComponentData().Values[0]

	key := fmt.Sprintf("thanks|%s|%s", i.GuildID, i.Member.User.ID)
	if err := c.CacheService.Get(context.Background(), key, &thanksData); err != nil {
		logger.Error("Failed to get cache", err)
		c.MessageService.SendStandardResponse(i, "Failed to get your choice, please try again", true, false)
		return
	}

	thanksData.CoreValue = coreValue

	if err := c.CacheService.Put(context.Background(), key, thanksData, 0); err != nil {
		logger.Error("Failed to put cache", err)
		c.MessageService.SendStandardResponse(i, "Failed to save your choice, please try again", true, false)
		return
	}

	reasonInput := discordgo.TextInput{
		CustomID:    "thanks_reason",
		Label:       "Reason",
		Style:       discordgo.TextInputParagraph,
		Placeholder: "Write your reason here",
		Required:    true,
	}

	modal := discordgo.InteractionResponseData{
		Title:    "Reason",
		CustomID: "thanks_modal_reason",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					reasonInput,
				},
			},
		},
	}

	c.MessageService.SendModalResponse(i, &modal)
}

func (c *ThanksCommand) sendThanksReason(i *discordgo.Interaction) {
	var thanksData ThanksData
	reason := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	key := fmt.Sprintf("thanks|%s|%s", i.GuildID, i.Member.User.ID)

	if err := c.CacheService.Get(context.Background(), key, &thanksData); err != nil {
		logger.Error("Failed to get cache", err)
		c.MessageService.SendStandardResponse(i, "something went wrong, please try again later", true, false)
		return
	}

	thanksData.Reason = reason

	if _, err := c.CacheService.Delete(context.Background(), key); err != nil {
		logger.Error("Failed to delete cache", err)
		c.MessageService.SendStandardResponse(i, "something went wrong, please try again later", true, false)
		return
	}

	c.proceedThanks(i, thanksData)
}

func (c *ThanksCommand) proceedThanks(i *discordgo.Interaction, thanksData ThanksData) {
	var response string
	ctx := context.Background()

	pointLogChannel, err := c.SettingService.GetPointLogChannel(ctx, i.GuildID)
	if err != nil {
		response = "Something went wrong, please try again later"
		logger.Error(response, err)
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	if pointLogChannel == "" {
		response = "Aktifkan fitur ini terlebih dahulu. Silahkan hubungi admin untuk aktifkan fitur ini."
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	if err := c.ThanksService.SendThanks(ctx, i.GuildID, i.Member.User.ID, thanksData.To, thanksData.CoreValue, thanksData.Reason); err != nil {
		if err == point.ErrLimitReached {
			response = "Limit mingguan kamu sudah habis, kamu bisa gunakan thanks lagi mulai senin depan"
		} else {
			response = "Something went wrong, can not add point"
			logger.Error(response, err)
		}
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	key := fmt.Sprintf("thanks_this_week|%s|%s|%s", i.GuildID, i.Member.User.ID, thanksData.To)
	if err := c.CacheService.Put(context.Background(), key, true, 0); err != nil {
		logger.Error("Failed to put cache", err)
		c.MessageService.SendStandardResponse(i, "Failed to save your choice, please try again", true, false)
		return
	}

	c.MessageService.SendStandardResponse(i, "Success", true, false)

	c.MessageService.SendStandardMessage(pointLogChannel, fmt.Sprintf("[%s] - <@%s> barusan kasih 10 rubic ke <@%s> karena %s", strings.ToUpper(thanksData.CoreValue), i.Member.User.ID, thanksData.To, thanksData.Reason))
	c.MessageService.SendStandardMessage(pointLogChannel, fmt.Sprintf("<@%s> barusan dapat 5 rubic karena sudah kirim thanks", i.Member.User.ID))

	if err := c.ThanksService.Log(ctx, i.GuildID, i.Member.User.ID, thanksData.To, thanksData.CoreValue, thanksData.Reason); err != nil {
		logger.Error("Failed to create thanks log", err)
	}
}
