package commands

import (
	"devteambot/config"
	"devteambot/internal/adapter/discord"
	"devteambot/internal/adapter/repository/redis"
	"devteambot/internal/application/service"
	"devteambot/internal/domain/message"
	"devteambot/internal/domain/point"
	"devteambot/internal/domain/review"
	"devteambot/internal/domain/setting"
	"devteambot/internal/pkg/cache"
	"devteambot/internal/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Conf     *config.Config `inject:"config"`
	Discord  *discord.App   `inject:"discord"`
	Cache    cache.Service  `inject:"cache"`
	RedisKey redis.RedisKey `inject:"redisKey"`
	cmdList  []*discordgo.ApplicationCommand

	MessageService message.Service   `inject:"messageService"`
	PointService   point.Service     `inject:"pointService"`
	ReviewService  review.Service    `inject:"reviewService"`
	SettingService setting.Service   `inject:"settingService"`
	AIService      service.AIService `inject:"aiService"`
}

func (c *Command) Startup() error {
	if c.Conf.Discord.RunInitCommand {
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
			{
				Name:        "thanks",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Kamu bisa kasih rubic dengan say thanks ke anggota tim yang lain",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "to",
						Description: "Orang yang dituju",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    true,
					},
					{
						Name:        "core",
						Description: "Pilih core value yang berhubungan",
						Type:        discordgo.ApplicationCommandOptionString,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Run",
								Value: "run",
							},
							{
								Name:  "Unity",
								Value: "unity",
							},
							{
								Name:  "Bravery",
								Value: "bravery",
							},
							{
								Name:  "Integrity",
								Value: "integrity",
							},
							{
								Name:  "Customer Oriented",
								Value: "customer-oriented",
							},
						},
						Required: true,
					},
					{
						Name:        "reason",
						Description: "Alasan kamu memberikan rubic",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
			{
				Name:        "thanks_leaderboard",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Show rubic leaderboard per category",
			},
			{
				Name:        "ask",
				Type:        discordgo.ChatApplicationCommand,
				Description: "Ask the bot a question",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "prompt",
						Description: "Your question",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
		}

		logger.Info("Adding commands...")
		for _, v := range commands {
			cmd, err := c.Discord.Bot.ApplicationCommandCreate(c.Conf.Discord.AppID, "", v)
			if err != nil {
				logger.Fatal(fmt.Sprintf("Cannot create command: %v %s", v.Name, err.Error()), err)
			}

			c.cmdList = append(c.cmdList, cmd)
		}
	}

	c.Discord.Bot.AddHandler(c.HandleCommand)

	return nil
}

func (c *Command) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "ping":
			c.Ping(s, i)
			return
		case "titip_review":
			c.TitipReview(s, i)
			return
		case "antrian_review":
			c.AntrianReview(s, i)
			return
		case "sudah_direview":
			c.SudahDireview(s, i)
			return
		case "thanks":
			c.Thanks(s, i)
			return
		case "thanks_leaderboard":
			c.ThanksLeaderboard(s, i)
			return
		}
	}

	if i.Type == discordgo.InteractionMessageComponent {
		commandID := i.MessageComponentData().CustomID
		switch {
		case strings.HasPrefix(commandID, "claim_role"):
			c.ClaimRole(s, i)
			return
		}
	}
}

func (c *Command) Shutdown() error {
	// if c.Conf.Discord.RunDeleteCommand {
	// 	for _, cmd := range c.cmdList {
	// 		err := c.Session.ApplicationCommandDelete(c.Conf.Discord.AppID, "", cmd.ID)
	// 		if err != nil {
	// 			logger.Fatal(fmt.Sprintf("Cannot delete command: %v %s", cmd.Name, err.Error()), err)
	// 		}
	// 	}
	// }

	return nil
}
