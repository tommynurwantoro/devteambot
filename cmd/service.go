package cmd

import (
	"devteambot/config"
	"devteambot/internal/bootstrap"
	"devteambot/internal/pkg/logger"
	"log"

	"github.com/spf13/cobra"
)

var (
	configFile string
)

var commandService = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Run service",
	Run: func(c *cobra.Command, args []string) {
		conf, err := config.Load(configFile)
		if err != nil {
			log.Fatal("Something went wrong when populating config")
		}

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
		bootstrap.Run(conf)
	},
}

func init() {
	commandService.Flags().StringVar(&configFile, "config", "./config.yaml", "Set config file path")
}

func GetCommand() *cobra.Command {
	return commandService
}
