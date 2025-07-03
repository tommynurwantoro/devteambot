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

type NewThanksCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Command    *Command `inject:"commandMember"`

	CacheService   cache.Service          `inject:"cache"`
	PointService   service.PointService   `inject:"pointService"`
	SettingService service.SettingService `inject:"settingService"`
	MessageService service.MessageService `inject:"messageService"`
}

func (c *NewThanksCommand) Startup() error {
	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *NewThanksCommand) Shutdown() error { return nil }

func (c *NewThanksCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

func (c *NewThanksCommand) sendThanks(i *discordgo.Interaction) {
	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.UserSelectMenu,
		CustomID:    "thanks_choose_user",
		Placeholder: "Select user",
		Options:     []discordgo.SelectMenuOption{},
	}

	c.MessageService.SendEmbedResponse(i, &discordgo.MessageSend{
		Content: "Please choose which user you want to send 🍪",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					selectMenu,
				},
			},
		},
	}, true)
}

func (c *NewThanksCommand) chooseUser(i *discordgo.Interaction) {
	to := i.MessageComponentData().Values[0]

	if to == i.Member.User.ID {
		c.MessageService.SendStandardResponse(i, "Tidak bisa berterima kasih ke diri sendiri", true, false)
		return
	}

	thanksData := ThanksData{
		To: to,
	}

	key := fmt.Sprintf("thanks_%s|%s", i.GuildID, i.Member.User.ID)
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

func (c *NewThanksCommand) chooseCoreValue(i *discordgo.Interaction) {
	var thanksData ThanksData
	coreValue := i.MessageComponentData().Values[0]

	key := fmt.Sprintf("thanks_%s|%s", i.GuildID, i.Member.User.ID)

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

func (c *NewThanksCommand) sendThanksReason(i *discordgo.Interaction) {
	var thanksData ThanksData
	reason := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	key := fmt.Sprintf("thanks_%s|%s", i.GuildID, i.Member.User.ID)

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

func (c *NewThanksCommand) proceedThanks(i *discordgo.Interaction, thanksData ThanksData) {
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
		response = "Aktifkan fitur ini terlebih dahulu dengan command /activate_point_feature"
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	if err := c.PointService.SendThanks(ctx, i.GuildID, i.Member.User.ID, thanksData.To, thanksData.CoreValue, thanksData.Reason); err != nil {
		if err == point.ErrLimitReached {
			response = "Limit mingguan kamu sudah habis, kamu bisa pakai command /thanks lagi mulai senin depan"
		} else {
			response = "Something went wrong, can not add rubic"
			logger.Error(response, err)
		}
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	c.MessageService.SendStandardResponse(i, "Success", true, false)

	c.MessageService.SendStandardMessage(pointLogChannel, fmt.Sprintf("[%s] - <@%s> barusan kasih 10 rubic ke <@%s> karena %s", strings.ToUpper(thanksData.CoreValue), i.Member.User.ID, thanksData.To, thanksData.Reason))
}
