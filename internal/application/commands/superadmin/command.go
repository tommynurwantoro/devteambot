package superadmin

import (
	"devteambot/config"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Conf    *config.Config `inject:"config"`
	Discord *discord.App   `inject:"discord"`
	cmdList []*discordgo.ApplicationCommand
}

func (c *Command) Startup() error {
	if c.Conf.Discord.RunInitCommand {
		c.RunAddCommand()
	}

	return nil
}

func (c *Command) Shutdown() error { return nil }

func (c *Command) AppendCommand(cmd *discordgo.ApplicationCommand) {
	serverManager := int64(discordgo.PermissionManageServer)
	dmPermission := new(bool)
	*dmPermission = false

	cmd.DefaultMemberPermissions = &serverManager
	cmd.DMPermission = dmPermission

	c.cmdList = append(c.cmdList, cmd)
}

func (c *Command) RunAddCommand() {
	logger.Info("Adding superadmin commands...")
	for _, v := range c.cmdList {
		logger.Info(fmt.Sprintf("Adding command: %v", v.Name))
		cmd, err := c.Discord.Bot.ApplicationCommandCreate(c.Conf.Discord.AppID, "", v)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Cannot create command: %v %s", v.Name, err.Error()), err)
		}

		c.cmdList = append(c.cmdList, cmd)
	}
}

// commands := []*discordgo.ApplicationCommand{
// 	{
// 		Name:                     "setup",
// 		Type:                     discordgo.ChatApplicationCommand,
// 		Description:              "Setup bot for the first time",
// 		DefaultMemberPermissions: &serverManager,
// 	},
// 	{
// 		Name:        "send_embed",
// 		Type:        discordgo.ChatApplicationCommand,
// 		Description: "Send embed message",
// 		Options: []*discordgo.ApplicationCommandOption{
// 			{
// 				Name:        "title",
// 				Description: "Message title",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    true,
// 			},
// 			{
// 				Name:        "content",
// 				Description: "Message content \"|\" for new line",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 			{
// 				Name:        "thumbnail",
// 				Description: "Thumbnail URL",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 		},
// 		DefaultMemberPermissions: &serverManager,
// 	},
// 	{
// 		Name:        "edit_embed",
// 		Description: "Edit embed message",
// 		Options: []*discordgo.ApplicationCommandOption{
// 			{
// 				Name:        "message_id",
// 				Description: "Message that will be edited",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    true,
// 			},
// 			{
// 				Name:        "title",
// 				Description: "Message title",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 			{
// 				Name:        "content",
// 				Description: "Message content \"|\" for new line",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 			{
// 				Name:        "thumbnail",
// 				Description: "Thumbnail URL",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 		},
// 		DefaultMemberPermissions: &serverManager,
// 	},
// 	{
// 		Name:        "add_button_feature",
// 		Type:        discordgo.ChatApplicationCommand,
// 		Description: "Add button feature to existing message",
// 		Options: []*discordgo.ApplicationCommandOption{
// 			{
// 				Name:        "command",
// 				Description: "Choose feature command",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Choices: []*discordgo.ApplicationCommandOptionChoice{
// 					{
// 						Name:  "Claim Role",
// 						Value: "claim_role",
// 					},
// 				},
// 				Required: true,
// 			},
// 			{
// 				Name:        "message_id",
// 				Description: "Message ID",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    true,
// 			},
// 			{
// 				Name:        "button1",
// 				Description: "RoleID|Name|IconName|IconID",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    true,
// 			},
// 			{
// 				Name:        "button2",
// 				Description: "RoleID|Name|IconName|IconID",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 			{
// 				Name:        "button3",
// 				Description: "RoleID|Name|IconName|IconID",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 			{
// 				Name:        "button4",
// 				Description: "RoleID|Name|IconName|IconID",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 			{
// 				Name:        "button5",
// 				Description: "RoleID|Name|IconName|IconID",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    false,
// 			},
// 		},
// 		DefaultMemberPermissions: &serverManager,
// 	},
// 	{
// 		Name:        "delete_button_feature",
// 		Type:        discordgo.ChatApplicationCommand,
// 		Description: "Delete button to existing message claim role",
// 		Options: []*discordgo.ApplicationCommandOption{
// 			{
// 				Name:        "message_id",
// 				Description: "Message ID",
// 				Type:        discordgo.ApplicationCommandOptionString,
// 				Required:    true,
// 			},
// 			{
// 				Name:        "index",
// 				Description: "First button is 1",
// 				Type:        discordgo.ApplicationCommandOptionInteger,
// 				Required:    true,
// 			},
// 		},
// 		DefaultMemberPermissions: &serverManager,
// 	},
// 	{
// 		Name:        "activate_reminder_sholat",
// 		Type:        discordgo.ChatApplicationCommand,
// 		Description: "Activate reminder sholat feature",
// 		Options: []*discordgo.ApplicationCommandOption{
// 			{
// 				Name:        "channel",
// 				Description: "Channel for the reminder",
// 				Type:        discordgo.ApplicationCommandOptionChannel,
// 				Required:    true,
// 			},
// 			{
// 				Name:        "role",
// 				Description: "Reminder sholat will mention this role",
// 				Type:        discordgo.ApplicationCommandOptionRole,
// 				Required:    true,
// 			},
// 		},
// 		DefaultMemberPermissions: &serverManager,
// 	},
// 	{
// 		Name:        "activate_reminder_presensi",
// 		Type:        discordgo.ChatApplicationCommand,
// 		Description: "Activate reminder presensi feature",
// 		Options: []*discordgo.ApplicationCommandOption{
// 			{
// 				Name:        "channel",
// 				Description: "Channel for the reminder",
// 				Type:        discordgo.ApplicationCommandOptionChannel,
// 				Required:    true,
// 			},
// 			{
// 				Name:        "role",
// 				Description: "Reminder presensi will mention this role",
// 				Type:        discordgo.ApplicationCommandOptionRole,
// 				Required:    true,
// 			},
// 		},
// 		DefaultMemberPermissions: &serverManager,
// 	},
// 	{
// 		Name:        "activate_point_feature",
// 		Type:        discordgo.ChatApplicationCommand,
// 		Description: "Activate point feature",
// 		Options: []*discordgo.ApplicationCommandOption{
// 			{
// 				Name:        "channel_id",
// 				Description: "Channel",
// 				Type:        discordgo.ApplicationCommandOptionChannel,
// 				Required:    true,
// 			},
// 		},
// 		DefaultMemberPermissions: &serverManager,
// 	},
// }
