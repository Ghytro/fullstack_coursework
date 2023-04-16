package middleware

import (
	"context"
	"errors"

	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type IUserProvider interface {
	Auth(ctx context.Context, username string, password string) (entity.PK, error)
}

func Auth(repo IUserProvider) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		jwtToken, ok := ctx.Locals("user_jwt").(*jwt.Token)
		if !ok {
			return entity.ErrRespUnauthorized(errors.New("unable to get jwt"))
		}
		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return entity.ErrRespUnauthorized(errors.New("unable to get claims from jwt"))
		}
		userName, ok := claims["user_name"].(string)
		if !ok {
			return entity.ErrRespUnauthorized(errors.New("unable to get 'id' from claims"))
		}
		userPass, ok := claims["password"].(string)
		if !ok {
			return entity.ErrRespUnauthorized(errors.New("unable to get 'pass' from claims"))
		}

		user, err := repo.Auth(ctx.Context(), userName, userPass)
		if err != nil {
			return entity.ErrRespUnauthorized(errors.New("incorrect token, auth again"))
		}
		ctx.Locals("user_entity", user)
		return ctx.Next()
	}
}
