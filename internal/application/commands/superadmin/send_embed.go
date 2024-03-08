package superadmin

import (
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/constant"
	"devteambot/internal/pkg/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type SendEmbedCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	MessageService service.MessageService `inject:"messageService"`
}

func (c *SendEmbedCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "send_embed",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Send embed message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "Message title",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "content",
				Description: "Message content \"|\" for new line",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "thumbnail",
				Description: "Thumbnail URL",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}

	c.CommandSuperAdmin.AppendCommand(c.AppCommand)
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *SendEmbedCommand) Shutdown() error { return nil }

func (c *SendEmbedCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *SendEmbedCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	var response string

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}

	var title, content, thumbnail string

	if opt, ok := optionMap["title"]; ok {
		title = opt.StringValue()
	}

	if opt, ok := optionMap["content"]; ok {
		content = opt.StringValue()
		content = strings.ReplaceAll(content, "|", "\n")
	}

	if opt, ok := optionMap["thumbnail"]; ok {
		thumbnail = opt.StringValue()
	}

	message := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       title,
				Color:       constant.BLUE,
				Description: content,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: thumbnail,
				},
			},
		},
	}

	_, err := s.ChannelMessageSendComplex(i.ChannelID, message)
	if err != nil {
		logger.Error(err.Error(), err)
		response = "Failed to send embed"
		c.MessageService.SendStandardResponse(i, response, true, false)
		return
	}

	response = "Success to send embed"
	c.MessageService.SendStandardResponse(i, response, true, false)
}
