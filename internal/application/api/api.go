package api

import "devteambot/internal/adapter/rest"

type Api struct {
	Rest *rest.Fiber `inject:"rest"`
}

func (a *Api) Startup() error {
	// v1 := a.Rest.Group("/v1")

	return nil
}

func (a *Api) Shutdown() error {
	return nil
}
