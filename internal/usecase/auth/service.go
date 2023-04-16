package auth

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/entity"

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
	if err := user.Validate(); err != nil {
		return "", err
	}
	_, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"user_name": user.Username,
		"password":  user.Password,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(s.jwtSecret)
	return token, err
}
