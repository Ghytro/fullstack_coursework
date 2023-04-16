package auth

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/entity"
)

type UseCase interface {
	MakeAuth(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, user *entity.User) (string, error)
}
