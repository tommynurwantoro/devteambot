package discord

import (
	"devteambot/config"
	"devteambot/internal/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type App struct {
	Bot    *discordgo.Session `inject:"botSession"`
	Config config.Discord     `inject:"discordConfig"`
}

func (a *App) Startup() error {
	a.Bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	err := a.Bot.Open()
	if err != nil {
		logger.Fatal("Error opening connection", err)
	}

	a.Bot.UpdateListeningStatus("Khuga")

	return nil
}

func (a *App) Shutdown() error {
	err := a.Bot.Close()
	if err != nil {
		return err
	}

	return nil
}
