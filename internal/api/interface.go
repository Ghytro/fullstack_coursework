package api

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	Routers(app fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler)
}
