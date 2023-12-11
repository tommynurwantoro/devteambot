package api

import "devteambot/internal/adapter/rest"

type Api struct {
	Rest          *rest.Fiber `inject:"rest"`
	SholatHandler SholatAPI   `inject:"sholatHandler"`
}

func (a *Api) Startup() error {
	v1 := a.Rest.Group("/v1")
	v1.Post("/get-sholat-schedule", a.SholatHandler.GetSholatSchedule)

	return nil
}

func (a *Api) Shutdown() error {
	return nil
}
