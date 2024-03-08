package bootstrap

import (
	"devteambot/config"
	"devteambot/internal/application/api"
	"devteambot/internal/application/commands"
	commandsuperadmin "devteambot/internal/application/commands/superadmin"
	"devteambot/internal/application/events"
	"devteambot/internal/application/scheduler"
	"devteambot/internal/application/service"
)

func RegisterScheduler(conf *config.Config) {
	appContainer.RegisterService("scheduler", new(scheduler.Scheduler))

	appContainer.RegisterService("presensiScheduler", new(scheduler.PresensiScheduler))
	appContainer.RegisterService("poinScheduler", new(scheduler.PoinScheduler))
	appContainer.RegisterService("sholatScheduler", new(scheduler.SholatScheduler))
}

func RegisterService() {
	appContainer.RegisterService("messageService", new(service.MessageService))
	appContainer.RegisterService("settingService", new(service.SettingService))
	appContainer.RegisterService("reviewService", new(service.ReviewService))
	appContainer.RegisterService("presensiService", new(service.PresensiService))
	appContainer.RegisterService("pointService", new(service.PointService))
	appContainer.RegisterService("sholatService", new(service.SholatService))
	appContainer.RegisterService("aiService", new(service.AI))
}

func RegisterAPI() {
	appContainer.RegisterService("api", new(api.Api))
}

func RegisterCommand() {
	appContainer.RegisterService("baseCommand", new(commands.Command))
	appContainer.RegisterService("commandSuperAdmin", new(commandsuperadmin.CommandSuperAdmin))
	appContainer.RegisterService("event", new(events.Event))
}
