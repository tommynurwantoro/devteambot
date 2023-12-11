package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devteambot/config"
	"devteambot/internal/adapter/rest"
	"devteambot/internal/adapter/scheduler"
	"devteambot/internal/application/commands"
	commandsuperadmin "devteambot/internal/application/commands/superadmin"
	"devteambot/internal/application/events"
	"devteambot/internal/constant"
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

	appContainer.RegisterService("redisKey", constant.NewRedisKey())
	appContainer.RegisterService("settingKey", constant.NewSettingKey())
	appContainer.RegisterService("color", constant.NewColor())

	// Dependency Injection
	RegisterDatabase()
	RegisterCache()
	RegisterDomain()
	RegisterDiscord()
	RegisterResty()
	RegisterRest()
	RegisterService()
	RegisterAPI()

	appContainer.RegisterService("baseCommand", new(commands.Command))
	appContainer.RegisterService("commandSuperAdmin", new(commandsuperadmin.CommandSuperAdmin))
	appContainer.RegisterService("event", new(events.Event))

	appContainer.RegisterService("scheduler", new(scheduler.Scheduler))

	// Check service readiness
	if err := appContainer.Ready(); err != nil {
		logger.Panic("Failed to populate service", err)
	}

	// Start server
	fiberApp := appContainer.GetServiceOrNil("rest").(*rest.Fiber)
	errs := make(chan error, 2)
	go func() {
		fmt.Printf("Listening on port :%d", conf.Http.Port)
		errs <- fiberApp.Listen(fmt.Sprintf(":%d", conf.Http.Port))
	}()

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

	logger.Info("Cleaning up resources...")

	appContainer.Shutdown()

	logger.Info("Bye")
}
