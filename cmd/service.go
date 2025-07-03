package cmd

import (
	"devteambot/config"
	"devteambot/internal/bootstrap"
	"devteambot/internal/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	runResetCommand bool
)

var commandService = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Run service",
	Run: func(c *cobra.Command, args []string) {
		godotenv.Load(".env")

		conf := config.Config{}
		conf.Load()
		conf.Discord.RunResetCommand = runResetCommand

		loggerConfig := logger.Config{
			App:           conf.AppName,
			AppVer:        conf.AppVersion,
			Env:           conf.Environment,
			FileLocation:  conf.Logger.FileLocation,
			FileMaxSize:   conf.Logger.FileMaxAge,
			FileMaxBackup: conf.Logger.FileMaxBackup,
			FileMaxAge:    conf.Logger.FileMaxAge,
			Stdout:        conf.Logger.Stdout,
		}

		logger.Load(loggerConfig)
		bootstrap.Run(&conf)
	},
}

func init() {
	commandService.Flags().BoolVar(&runResetCommand, "reset-command", false, "Run reset command")
}

func GetCommand() *cobra.Command {
	return commandService
}
