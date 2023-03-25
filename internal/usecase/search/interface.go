package search

import (
	"context"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/repository"
)

type UserRepository interface {
	GetUserListSearch(ctx context.Context, filter *repository.UserSearchFilter) ([]*entity.User, error)

	RunInTransaction(ctx context.Context, fn func(tx *database.TX) error) error
	WithTX(tx *database.TX) *repository.UserRepository
}

type PollsRepository interface {
	GetPollListSearch(ctx context.Context, filter *repository.PollSearchFilter) ([]*entity.Poll, error)

	RunInTransaction(ctx context.Context, fn func(tx *database.TX) error) error
	WithTX(tx *database.TX) *repository.PollsRepo
}
