package member

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

// {
// 	Name:        "ping",
// 	Type:        discordgo.ChatApplicationCommand,
// 	Description: "Check bot server status",
// },
// {
// 	Name:        "titip_review",
// 	Type:        discordgo.ChatApplicationCommand,
// 	Description: "Titip review ke anggota squad",
// 	Options: []*discordgo.ApplicationCommandOption{
// 		{
// 			Name:        "title",
// 			Description: "Title",
// 			Type:        discordgo.ApplicationCommandOptionString,
// 			Required:    true,
// 		},
// 		{
// 			Name:        "url",
// 			Description: "Url",
// 			Type:        discordgo.ApplicationCommandOptionString,
// 			Required:    true,
// 		},
// 		{
// 			Name:        "reviewer_1",
// 			Description: "Reviewer 1",
// 			Type:        discordgo.ApplicationCommandOptionUser,
// 			Required:    true,
// 		},
// 		{
// 			Name:        "reviewer_2",
// 			Description: "Reviewer 2",
// 			Type:        discordgo.ApplicationCommandOptionUser,
// 			Required:    false,
// 		},
// 		{
// 			Name:        "reviewer_3",
// 			Description: "Reviewer 3",
// 			Type:        discordgo.ApplicationCommandOptionUser,
// 			Required:    false,
// 		},
// 		{
// 			Name:        "reviewer_4",
// 			Description: "Reviewer 4",
// 			Type:        discordgo.ApplicationCommandOptionUser,
// 			Required:    false,
// 		},
// 		{
// 			Name:        "reviewer_5",
// 			Description: "Reviewer 5",
// 			Type:        discordgo.ApplicationCommandOptionUser,
// 			Required:    false,
// 		},
// 	},
// },
// {
// 	Name:        "antrian_review",
// 	Type:        discordgo.ChatApplicationCommand,
// 	Description: "Show all pending review",
// },
// {
// 	Name:        "sudah_direview",
// 	Type:        discordgo.ChatApplicationCommand,
// 	Description: "Sudah melakukan review",
// 	Options: []*discordgo.ApplicationCommandOption{
// 		{
// 			Name:        "number",
// 			Description: "Antrian nomor berapa",
// 			Type:        discordgo.ApplicationCommandOptionInteger,
// 			Required:    true,
// 		},
// 	},
// },
// {
// 	Name:        "thanks",
// 	Type:        discordgo.ChatApplicationCommand,
// 	Description: "Kamu bisa kasih rubic dengan say thanks ke anggota tim yang lain",
// 	Options: []*discordgo.ApplicationCommandOption{
// 		{
// 			Name:        "to",
// 			Description: "Orang yang dituju",
// 			Type:        discordgo.ApplicationCommandOptionUser,
// 			Required:    true,
// 		},
// 		{
// 			Name:        "core",
// 			Description: "Pilih core value yang berhubungan",
// 			Type:        discordgo.ApplicationCommandOptionString,
// 			Choices: []*discordgo.ApplicationCommandOptionChoice{
// 				{
// 					Name:  "Run",
// 					Value: "run",
// 				},
// 				{
// 					Name:  "Unity",
// 					Value: "unity",
// 				},
// 				{
// 					Name:  "Bravery",
// 					Value: "bravery",
// 				},
// 				{
// 					Name:  "Integrity",
// 					Value: "integrity",
// 				},
// 				{
// 					Name:  "Customer Oriented",
// 					Value: "customer-oriented",
// 				},
// 			},
// 			Required: true,
// 		},
// 		{
// 			Name:        "reason",
// 			Description: "Alasan kamu memberikan rubic",
// 			Type:        discordgo.ApplicationCommandOptionString,
// 			Required:    true,
// 		},
// 	},
// },
// {
// 	Name:        "thanks_leaderboard",
// 	Type:        discordgo.ChatApplicationCommand,
// 	Description: "Show rubic leaderboard per category",
// },
// {
// 	Name:        "ask",
// 	Type:        discordgo.ChatApplicationCommand,
// 	Description: "Ask the bot a question",
// 	Options: []*discordgo.ApplicationCommandOption{
// 		{
// 			Name:        "prompt",
// 			Description: "Your question",
// 			Type:        discordgo.ApplicationCommandOptionString,
// 			Required:    true,
// 		},
// 	},
// },
