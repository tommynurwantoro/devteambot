package superadmin

import (
	"devteambot/internal/application/commands"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CommandSuperAdmin struct {
	Command *commands.Command `inject:"baseCommand"`
	cmdList []*discordgo.ApplicationCommand

	SettingRepository setting.Repository `inject:"settingRepository"`
	SettingService    setting.Service    `inject:"settingService"`
}

func (c *CommandSuperAdmin) Startup() error {
	serverManager := int64(discordgo.PermissionManageServer)

	if c.Command.Conf.Discord.RunInitCommand {
		commands := []*discordgo.ApplicationCommand{
			{
				Name:                     "setup",
				Type:                     discordgo.ChatApplicationCommand,
				Description:              "Setup bot for the first time",
				DefaultMemberPermissions: &serverManager,
			},
			{
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
				DefaultMemberPermissions: &serverManager,
			},
			{
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
				DefaultMemberPermissions: &serverManager,
			},
			{
				Name:        "add_button_feature",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Add button feature to existing message",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "command",
						Description: "Choose feature command",
						Type:        discordgo.ApplicationCommandOptionString,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Claim Role",
								Value: "claim_role",
							},
						},
						Required: true,
					},
					{
						Name:        "message_id",
						Description: "Message ID",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "button1",
						Description: "RoleID|Name|IconName|IconID",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "button2",
						Description: "RoleID|Name|IconName|IconID",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    false,
					},
					{
						Name:        "button3",
						Description: "RoleID|Name|IconName|IconID",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    false,
					},
					{
						Name:        "button4",
						Description: "RoleID|Name|IconName|IconID",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    false,
					},
					{
						Name:        "button5",
						Description: "RoleID|Name|IconName|IconID",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    false,
					},
				},
				DefaultMemberPermissions: &serverManager,
			},
			{
				Name:        "delete_button_feature",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Delete button to existing message claim role",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "message_id",
						Description: "Message ID",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "index",
						Description: "First button is 1",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Required:    true,
					},
				},
				DefaultMemberPermissions: &serverManager,
			},
			{
				Name:        "activate_reminder_sholat",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Activate reminder sholat feature",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "channel",
						Description: "Channel for the reminder",
						Type:        discordgo.ApplicationCommandOptionChannel,
						Required:    true,
					},
					{
						Name:        "role",
						Description: "Reminder sholat will mention this role",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
				DefaultMemberPermissions: &serverManager,
			},
			{
				Name:        "activate_reminder_presensi",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Activate reminder presensi feature",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "channel",
						Description: "Channel for the reminder",
						Type:        discordgo.ApplicationCommandOptionChannel,
						Required:    true,
					},
					{
						Name:        "role",
						Description: "Reminder presensi will mention this role",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
				DefaultMemberPermissions: &serverManager,
			},
			{
				Name:        "activate_point_feature",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Activate point feature",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "channel_id",
						Description: "Channel",
						Type:        discordgo.ApplicationCommandOptionChannel,
						Required:    true,
					},
				},
				DefaultMemberPermissions: &serverManager,
			},
		}

		logger.Info("Adding super admin commands...")
		for _, v := range commands {
			cmd, err := c.Command.Discord.Bot.ApplicationCommandCreate(c.Command.Conf.Discord.AppID, "", v)
			if err != nil {
				logger.Fatal(fmt.Sprintf("Cannot create command: %v %s", v.Name, err.Error()), err)
			}

			c.cmdList = append(c.cmdList, cmd)
		}
	}

	c.Command.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *CommandSuperAdmin) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "setup":
			c.Setup(s, i)
			return
		case "add_button_feature":
			c.AddButtonFeature(s, i)
			return
		case "delete_button_feature":
			c.DeleteButtonFeature(s, i)
			return
		case "edit_embed":
			c.EditEmbed(s, i)
			return
		case "send_embed":
			c.SendEmbed(s, i)
			return
		case "activate_reminder_sholat":
			c.ActivateReminderSholat(s, i)
			return
		case "activate_reminder_presensi":
			c.ActivateReminderPresensi(s, i)
			return
		case "activate_point_feature":
			c.ActivatePoint(s, i)
			return
		}
	}

	if i.Type == discordgo.InteractionModalSubmit {
		switch i.ModalSubmitData().CustomID {
		case "setup_superadmin":
			c.SetupSuperAdmin(s, i)
			return
		}
	}
}

func (c *CommandSuperAdmin) Shutdown() error {
	// if c.Conf.RunDeleteCommand {
	// for _, cmd := range c.cmdList {
	// 	err := c.Session.ApplicationCommandDelete(c.Conf.AppID, c.Conf.GuildID, cmd.ID)
	// 	if err != nil {
	// 		logger.Fatal(fmt.Sprintf("Cannot delete command: %v %s", cmd.Name, err.Error()), err)
	// 	}
	// }
	// }

	return nil
}
