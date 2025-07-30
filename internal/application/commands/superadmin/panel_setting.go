package superadmin

import (
	"devteambot/internal/application/service"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type PanelSettingCommand struct {
	AppCommand        *discordgo.ApplicationCommand
	CommandSuperAdmin *Command `inject:"commandSuperAdmin"`

	MessageService service.MessageService `inject:"messageService"`
	SettingService service.SettingService `inject:"settingService"`
}

func (c *PanelSettingCommand) Startup() error {
	c.AppCommand = &discordgo.ApplicationCommand{
		Name:        "panel_setting",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Open admin panel setting",
		Options:     []*discordgo.ApplicationCommandOption{},
	}

	c.CommandSuperAdmin.AppendCommand(c.AppCommand)
	c.CommandSuperAdmin.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *PanelSettingCommand) Shutdown() error { return nil }

func (c *PanelSettingCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand && i.ApplicationCommandData().Name == c.AppCommand.Name {
		c.Do(s, i.Interaction)
	}
}

func (c *PanelSettingCommand) Do(s *discordgo.Session, i *discordgo.Interaction) {
	embedMessage := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Admin Panel Setting",
			Description: "This is the admin panel setting",
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "🍪",
						},
						Label:    "Activate Point Feature",
						Style:    discordgo.PrimaryButton,
						CustomID: "activate_point_feature",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "✅",
						},
						Label:    "Activate Reminder Presensi Feature",
						Style:    discordgo.PrimaryButton,
						CustomID: "activate_reminder_presensi_feature",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "🕌",
						},
						Label:    "Activate Reminder Sholat Feature",
						Style:    discordgo.PrimaryButton,
						CustomID: "activate_reminder_sholat_feature",
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "💰",
						},
						Label:    "Marketplace New Item",
						Style:    discordgo.PrimaryButton,
						CustomID: "marketplace_new_item",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "💰",
						},
						Label:    "Marketplace Update Item",
						Style:    discordgo.PrimaryButton,
						CustomID: "marketplace_update_item",
					},
				},
			},
		},
	}

	if err := c.MessageService.SendEmbedMessage(i.ChannelID, embedMessage); err != nil {
		logger.Error("Failed to send embed message", err)
	}

	c.MessageService.SendStandardResponse(i, "Success to open admin panel setting", true, false)
}
