package main

import (
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/repository"
	"github.com/Ghytro/galleryapp/internal/view"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func serve(postgresUrl, mongoUrl string) {
	db := database.NewPGDB(postgresUrl, &database.PGLogger{})
	jwtSecret := []byte("")

	repo := repository.NewRepository(db, logrus.StandardLogger())
	NewApp(
		jwtSecret,
		db,
	).Listen(":3001")
}

func NewApp(token interface{}, db repository.DBI, views ...view.View) *fiber.App {
	r := fiber.New()
	authHandler := jwtware.New(jwtware.Config{
		SigningKey:     token,
		TokenLookup:    "cookie:jwt",
		ContextKey:     "user_jwt",
		SuccessHandler: authSuccessHandler(db),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/auth", fiber.StatusSeeOther)
		},
	})
	middlewares := []fiber.Handler{
		func(c *fiber.Ctx) error {
			c.Set("Content-Type", "text/html;charset=utf-8")
			return c.Next()
		},
	}
	for _, v := range views {
		v.Routers(r, authHandler, middlewares...)
	}
	return r
}
