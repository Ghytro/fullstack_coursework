package auth

import (
	"github.com/Ghytro/galleryapp/internal/api"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	service UseCase
}

func NewAPI(service UseCase) *API {
	return &API{
		service: service,
	}
}

func (a *API) Routers(app fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler) {
	r := fiber.New()
	for _, m := range middlewares {
		r.Use(m)
	}
	r.Post("/", a.register)
	r.Get("/", a.auth)
	app.Mount("/auth", r)
}

func (a *API) register(ctx *fiber.Ctx) error {
	return api.ErrNotImplemented
}

func (a *API) auth(ctx *fiber.Ctx) error {
	return api.ErrNotImplemented
}
