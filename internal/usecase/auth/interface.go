package auth

import (
	"context"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/repository"
)

type Repository interface {
	Auth(ctx context.Context, username string, password string) (entity.PK, error)
	CreateUser(ctx context.Context, user *entity.User) (entity.PK, error)
	GetUser(ctx context.Context, userID entity.PK) (*entity.User, error)

	RunInTransaction(ctx context.Context, f func(*database.TX) error) error
	WithTX(*database.TX) *repository.UserRepository
}
