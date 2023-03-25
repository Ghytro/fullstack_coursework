package auth

import (
	"errors"
	"fmt"
	"github.com/Ghytro/galleryapp/internal/entity"
	"os"

	"github.com/gofiber/fiber/v2"
)

type View struct {
	service UseCase
}

func NewView(s UseCase) *View {
	return &View{
		service: s,
	}
}

func (v *View) Routers(app fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler) {
	r := fiber.New()
	for _, m := range middlewares {
		r.Use(m)
	}
	r.Get("/", v.getAuth)
	r.Post("/", v.makeAuth)
	r.Patch("/", v.patchAuth)
	app.Mount("/auth", r)

	r = fiber.New()
	r.Post("/", v.register)
	app.Mount("/register", r)
}

func (v *View) getAuth(c *fiber.Ctx) error {
	file, err := os.Open("./web/auth/index.html")
	if err != nil {
		return entity.ErrRespInternalServerError(err)
	}
	return c.SendStream(file)
}

func (v *View) makeAuth(c *fiber.Ctx) error {
	fmt.Println("make auth handler")
	var model MakeAuthRequest
	if err := c.BodyParser(&model); err != nil {
		return entity.ErrRespIncorrectForm()
	}
	fmt.Println(model)
	token, err := v.service.MakeAuth(c.Context(), model.Username, model.Password)
	if err != nil {
		fmt.Println(err)
		return entity.ErrRespBadRequest(err)
	}
	fmt.Println(token)
	c.Cookie(&fiber.Cookie{Name: "jwt", Value: token})
	return c.Redirect("/profile", fiber.StatusSeeOther)
}

func (v *View) patchAuth(c *fiber.Ctx) error {
	return errors.New("missing impl") // TODO
}

func (v *View) register(c *fiber.Ctx) error {
	var model entity.User
	if err := c.BodyParser(&model); err != nil {
		return entity.ErrRespIncorrectForm()
	}

	token, err := v.service.Register(c.Context(), &model)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	c.Cookie(&fiber.Cookie{Name: "jwt", Value: token})
	return c.Redirect("/profile", fiber.StatusSeeOther)
}
