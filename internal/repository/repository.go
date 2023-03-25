package repository

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/database"
)

type MonoRepo struct {
	*UserRepository
}

func NewRepository(db DBI) IRepository {
	return &MonoRepo{
		UserRepository: NewUserRepo(db),
	}
}

func (r *MonoRepo) RunInTransaction(ctx context.Context, fn func(repo IRepository) error) error {
	return r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		txRepo := NewRepository(tx)
		return fn(txRepo)
	})
}

type IRepository interface {
	IUserRepository
}
