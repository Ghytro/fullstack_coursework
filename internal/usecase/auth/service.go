package auth

import (
	"context"
	"errors"

	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/validation"
	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	repo      Repository
	jwtSecret interface{}
}

func NewService(r Repository, secret interface{}) *Service {
	return &Service{
		repo:      r,
		jwtSecret: secret,
	}
}

func (s *Service) MakeAuth(ctx context.Context, username string, password string) (string, error) {
	userID, err := s.repo.Auth(ctx, username, password)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"id":   userID,
		"pass": password,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := t.SignedString(s.jwtSecret)
	return accessToken, err
}

func (s *Service) PatchAuth(ctx context.Context, username string, password string) (string, error) {
	return "missing impl", nil // TODO
}

func (s *Service) Register(ctx context.Context, user *entity.User) (string, error) {
	if err := validation.ValidateUser(user); err != nil {
		return "", err
	}
	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"id":   userID,
		"pass": user.Password,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(s.jwtSecret)
	return token, err
}

func (s *Service) AuthMiddleware(ctx *fiber.Ctx) error {
	var u entity.User
	jwtToken, ok := ctx.Locals("user_jwt").(*jwt.Token)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("unable to get jwt"))
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("unable to get claims from jwt"))
	}
	userIdClaims, ok := claims["id"].(float64)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("unable to get 'id' from claims"))
	}
	userId := entity.PK(userIdClaims)
	userPass, ok := claims["pass"].(string)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("unable to get 'pass' from claims"))
	}
	if err := db.ModelContext(c.Context(), &u).
		Where("id = ? AND password = crypt(?, password)", userId, userPass).
		Select(); err != nil {
		if err == pg.ErrNoRows {
			return entity.ErrRespUnauthorized(errors.New("incorrect token, auth again"))
		}
		return entity.ErrRespInternalServerError(err)
	}
	c.Locals("user_entity", &u)
	return c.Next()
}
