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

	if err := a.Bot.UpdateListeningStatus("RUNers"); err != nil {
		logger.Fatal("Failed to update listening status", err)
	}

	if a.Conf.Discord.RunResetCommand {
		if err := a.resetCommand(); err != nil {
			logger.Fatal("Failed to reset command", err)
		}
	}

	return nil
}

func (a *App) Shutdown() error {
	return a.Bot.Close()
}

func (a *App) resetCommand() error {
	logger.Info("Reset all commands...")

	logger.Info("Getting command list...")
	cmdList, err := a.Bot.ApplicationCommands(a.Conf.Discord.AppID, "")
	if err != nil {
		logger.Error("Cannot get command list", err)
		return err
	}

	logger.Info("Deleting commands...")
	for _, cmd := range cmdList {
		err := a.Bot.ApplicationCommandDelete(a.Conf.Discord.AppID, "", cmd.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("Cannot delete command: %v %s", cmd.Name, err.Error()), err)
			return err
		}
	}

	return nil
}
