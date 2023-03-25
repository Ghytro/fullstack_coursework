package view

import "github.com/gofiber/fiber/v2"

type View interface {
	Routers(app fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler)
}
