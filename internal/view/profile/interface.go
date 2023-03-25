package profile

import (
	"context"
	"github.com/Ghytro/galleryapp/internal/entity"
)

type UseCase interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.PK, error)
	GetUser(ctx context.Context, userID entity.PK) (*entity.User, error)
	GetUserWithPolls(ctx context.Context, userID entity.PK, limit int) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID entity.PK) error
}
