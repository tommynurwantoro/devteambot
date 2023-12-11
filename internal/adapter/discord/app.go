package discord

import (
	"devteambot/config"
	"devteambot/internal/pkg/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type App struct {
	Bot  *discordgo.Session
	Conf *config.Config `inject:"config"`
}

func (a *App) Startup() error {
	bot, err := discordgo.New("Bot " + a.Conf.Discord.Token)
	if err != nil {
		logger.Panic("Cannot authorizing bot ", err)
	}
	logger.Info("Bot initialized")

	a.Bot = bot
	a.Bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	err = a.Bot.Open()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error opening connection: %s", err.Error()), err)
	}

	a.Bot.UpdateListeningStatus("RUNers")

	return nil
}

func (a *App) Shutdown() error {
	return a.Bot.Close()
}
