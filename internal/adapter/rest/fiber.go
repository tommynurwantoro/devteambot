package rest

import (
	"devteambot/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Fiber struct {
	*fiber.App
	Conf *config.Config `inject:"config"`
}

func (f *Fiber) Startup() error {
	//starting http
	f.App = fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(f.Conf.Http.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(f.Conf.Http.WriteTimeout) * time.Second,
	})

	// Middleware
	f.App.Use(recover.New())
	f.App.Use(requestid.New())

	return nil
}

func (f *Fiber) Shutdown() error { return f.App.Shutdown() }
