package commands

import (
	"devteambot/config"
	"devteambot/internal/adapter/cache"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/constant"
	"devteambot/internal/domain/review"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
)

type Command struct {
	Conf        config.Discord      `inject:"discordConfig"`
	App         *discord.App        `inject:"discordApp"`
	Cache       cache.Cache         `inject:"cache"`
	SettingKey  constant.SettingKey `inject:"settingKey"`
	RedisKey    constant.RedisKey   `inject:"redisKey"`
	Color       constant.Color      `inject:"color"`
	GoogleSheet *resty.Client       `inject:"googleSheet"`
	Admins      map[string]bool     `inject:"admins"`
	SuperAdmins map[string]bool     `inject:"superAdmins"`
	cmdList     []*discordgo.ApplicationCommand

	SettingRepository setting.Repository `inject:"settingRepository"`
	ReviewRepository  review.Repository  `inject:"reviewRepository"`
}

func (c *Command) Startup() error {
	if c.Conf.RunInitCommand {
		commands := []*discordgo.ApplicationCommand{
			{
				Name:        "ping",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Check bot server status",
			},
			{
				Name:        "titip_review",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Titip review ke anggota squad",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "title",
						Description: "Title",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "url",
						Description: "Url",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "reviewer_1",
						Description: "Reviewer 1",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    true,
					},
					{
						Name:        "reviewer_2",
						Description: "Reviewer 2",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    false,
					},
					{
						Name:        "reviewer_3",
						Description: "Reviewer 3",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    false,
					},
					{
						Name:        "reviewer_4",
						Description: "Reviewer 4",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    false,
					},
					{
						Name:        "reviewer_5",
						Description: "Reviewer 5",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    false,
					},
				},
			},
			{
				Name:        "antrian_review",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Show all pending review",
			},
			{
				Name:        "sudah_direview",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Sudah melakukan review",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "number",
						Description: "Antrian nomor berapa",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Required:    true,
					},
				},
			},
		}

		logger.Info("Adding commands...")
		for _, v := range commands {
			cmd, err := c.App.Bot.ApplicationCommandCreate(c.Conf.AppID, "", v)
			if err != nil {
				logger.Fatal(fmt.Sprintf("Cannot create command: %v %s", v.Name, err.Error()), err)
			}

			c.cmdList = append(c.cmdList, cmd)
		}
	}

	c.App.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *Command) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "ping":
			c.Ping(s, i)
			break
		case "titip_review":
			c.TitipReview(s, i)
			break
		case "antrian_review":
			c.AntrianReview(s, i)
			break
		case "sudah_direview":
			c.SudahDireview(s, i)
			break
		}
	}

	if i.Type == discordgo.InteractionMessageComponent {
		commandID := i.MessageComponentData().CustomID
		switch {
		case strings.HasPrefix(commandID, "claim_role"):
			c.ClaimRole(s, i)
			break
		}
	}
}

func (c *Command) Shutdown() error {
	// if c.Conf.RunDeleteCommand {
	// 	for _, cmd := range c.cmdList {
	// 		err := c.Session.ApplicationCommandDelete(c.Conf.AppID, "", cmd.ID)
	// 		if err != nil {
	// 			logger.Fatal(fmt.Sprintf("Cannot delete command: %v %s", cmd.Name, err.Error()), err)
	// 		}
	// 	}
	// }

	return nil
}
