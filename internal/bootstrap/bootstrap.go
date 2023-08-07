package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devteambot/config"
	"devteambot/internal/adapter/scheduler"
	"devteambot/internal/commands"
	commandsuperadmin "devteambot/internal/commands/superadmin"
	"devteambot/internal/constant"
	"devteambot/internal/events"
	"devteambot/internal/pkg/container"
	"devteambot/internal/pkg/logger"

	"github.com/asaskevich/govalidator"
)

var appContainer = container.New()

func Run(conf *config.Config) {
	_, err := govalidator.ValidateStruct(conf)
	if err != nil {
		logger.Panic("invalid config", err)
	}

	logger.Info("Serving...")

	appContainer.RegisterService("config", conf)

	superAdmins := make(map[string]bool)
	admins := make(map[string]bool)
	// for _, id := range conf.Discord.SuperAdminRoleIDs {
	// 	superAdmins[id] = true
	// }
	// for _, id := range conf.Discord.AdminRoleIDs {
	// 	admins[id] = true
	// }
	appContainer.RegisterService("superAdmins", superAdmins)
	appContainer.RegisterService("admins", admins)
	appContainer.RegisterService("discordConfig", conf.Discord)

	appContainer.RegisterService("redisKey", constant.NewRedisKey())
	appContainer.RegisterService("settingKey", constant.NewSettingKey())
	appContainer.RegisterService("phrase", constant.NewPhrase())
	appContainer.RegisterService("color", constant.NewColor())

	// Dependency Injection
	RegisterDatabase(&conf.Database)
	RegisterCache(&conf.Redis)
	RegisterDomain()
	RegisterDiscord(&conf.Discord)
	// RegisterRest(conf)
	RegisterAPI()

	appContainer.RegisterService("baseCommand", new(commands.Command))
	appContainer.RegisterService("commandSuperAdmin", new(commandsuperadmin.CommandSuperAdmin))
	appContainer.RegisterService("event", new(events.Event))

	appContainer.RegisterService("scheduler", new(scheduler.Scheduler))

	// Check service readiness
	if err := appContainer.Ready(); err != nil {
		logger.Panic("Failed to populate service", err)
	}

	logger.Info(fmt.Sprintf("%s started", conf.AppName))

	GracefulShutdown(conf)
}

func GracefulShutdown(conf *config.Config) {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	delay := conf.ShutdownDelay

	logger.Info(fmt.Sprintf("Signal termination received. Waiting %v to shutdown.", delay))

	time.Sleep(delay)

	logger.Info(fmt.Sprintf("Cleaning up resources..."))

	appContainer.Shutdown()

	logger.Info("Bye")
}
