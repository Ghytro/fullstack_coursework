package repository

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/sirupsen/logrus"
)

type MonoRepo struct {
	db  *database.PGDB
	log *logrus.Logger
	*UserRepository
}

func NewRepository(db *database.PGDB, log *logrus.Logger) IRepository {
	r := &MonoRepo{
		db: db,
	}
	r.UserRepository = &UserRepository{r}
	return r
}

func (r *MonoRepo) RunInTransaction(ctx context.Context, fn func(ctx context.Context, repo IRepository) error) error {
	if tx := database.GetTx(ctx); tx == nil {
		tx, err := r.db.BeginContext(ctx)
		if err != nil {
			return err
		}
		ctx := database.ApplyTx(ctx, tx)
		if err := r.db.RunInTransaction(ctx, func(tx *database.TX) error {
			return fn(ctx, r)
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				r.log.Error(err)
			}
			return nil
		}
		if err := tx.Commit(); err != nil {
			r.log.Error(err)
		}
		return nil
	}
	return r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		return fn(ctx, r)
	})
}

type IRepository interface {
	IUserRepository
}
