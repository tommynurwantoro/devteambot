package superadmin

import (
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type EditEmbedCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	MessageService service.MessageService `inject:"messageService"`
}

func (c *EditEmbedCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "edit_embed",
		Description: "Edit embed message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "message_id",
				Description: "Message that will be edited",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "title",
				Description: "Message title",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
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

func (c *EditEmbedCommand) Shutdown() error { return nil }

func (c *EditEmbedCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *EditEmbedCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	var response string

	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name != "edit_embed" {
			return
		}

		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, option := range options {
			optionMap[option.Name] = option
		}

		var messageID, title, content, thumbnail string

		if opt, ok := optionMap["message_id"]; ok {
			messageID = opt.StringValue()
		}

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

		existingMessage, err := s.ChannelMessage(i.ChannelID, messageID)
		if err != nil {
			logger.Error(err.Error(), err)
			response = "Something went wrong, please try again later"
			c.MessageService.SendStandardResponse(i, response, true, false)
			return
		}

		embedMessage := existingMessage.Embeds[0]

		if title != "" {
			embedMessage.Title = title
		}

		if content != "" {
			embedMessage.Description = content
		}

		if thumbnail != "" {
			embedMessage.Thumbnail.URL = thumbnail
		}

		existingMessage.Embeds = []*discordgo.MessageEmbed{
			embedMessage,
		}

		editMessage := &discordgo.MessageEdit{
			Content:    &existingMessage.Content,
			Components: existingMessage.Components,
			Embeds:     existingMessage.Embeds,
			ID:         existingMessage.ID,
			Channel:    existingMessage.ChannelID,
		}

		_, err = s.ChannelMessageEditComplex(editMessage)
		if err != nil {
			logger.Error(err.Error(), err)
			response = "Failed to edit embed"
			c.MessageService.SendStandardResponse(i, response, true, false)
			return
		}

		response = "Success to edit embed"
		c.MessageService.SendStandardResponse(i, response, true, false)
	}
}
