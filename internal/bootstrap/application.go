package bootstrap

import (
	"devteambot/internal/application/api"
	"devteambot/internal/application/commands"
	commandsuperadmin "devteambot/internal/application/commands/superadmin"
	"devteambot/internal/application/events"
	"devteambot/internal/application/service"
)

func RegisterService() {
	appContainer.RegisterService("sholatService", new(service.SholatService))
}

func RegisterAPI() {
	appContainer.RegisterService("api", new(api.Api))
	appContainer.RegisterService("sholatHandler", new(api.SholatHandler))
}

func RegisterCommand() {
	appContainer.RegisterService("baseCommand", new(commands.Command))
	appContainer.RegisterService("commandSuperAdmin", new(commandsuperadmin.CommandSuperAdmin))
	appContainer.RegisterService("event", new(events.Event))
}
