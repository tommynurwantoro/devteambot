package cmd

import (
	"devteambot/config"
	"devteambot/internal/bootstrap"
	"devteambot/internal/pkg/conf"
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
		x := conf.NewConfig(configFile, new(config.Config))
		newConfig, ok := x.(*config.Config)
		if !ok {
			log.Fatal("Something went wrong when populating config")
		}
		logger.Load(newConfig.Environment)
		bootstrap.Run(newConfig)
	},
}

func init() {
	commandService.Flags().StringVar(&configFile, "config", "./config.yaml", "Set config file path")
}

func GetCommand() *cobra.Command {
	return commandService
}
