package main

import (
	"github.com/Ghytro/galleryapp/internal/api"
	apiAuth "github.com/Ghytro/galleryapp/internal/api/auth"
	"github.com/Ghytro/galleryapp/internal/api/middleware"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/repository"
	serviceAuth "github.com/Ghytro/galleryapp/internal/usecase/auth"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func serve(postgresUrl, mongoUrl string) {
	db := database.NewPGDB(postgresUrl, &database.PGLogger{})
	jwtSecret := []byte("")

	repo := repository.NewRepository(db, logrus.StandardLogger())
	authUseCase := serviceAuth.NewService(repo, jwtSecret)
	authApi := apiAuth.NewAPI(authUseCase)
	NewApp(
		jwtSecret,
		repo,
		authApi,
	).Listen(":3001")
}

func NewApp(token interface{}, repo middleware.IUserProvider, apis ...api.Handlers) *fiber.App {
	r := fiber.New()
	authHandler := jwtware.New(jwtware.Config{
		SigningKey:     token,
		TokenLookup:    "cookie:jwt",
		ContextKey:     "user_jwt",
		SuccessHandler: middleware.Auth(repo),
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
	for _, a := range apis {
		a.Routers(r, authHandler, middlewares...)
	}
	return r
}
