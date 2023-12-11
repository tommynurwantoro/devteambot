package bootstrap

import (
	"devteambot/internal/application/api"
	"devteambot/internal/application/service"
)

func RegisterService() {
	appContainer.RegisterService("sholatService", new(service.SholatService))
}

func RegisterAPI() {
	appContainer.RegisterService("api", new(api.Api))
	appContainer.RegisterService("sholatHandler", new(api.SholatHandler))
}
