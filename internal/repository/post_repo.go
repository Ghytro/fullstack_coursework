package repository

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/entity"
)

type PostRepository struct {
	*MonoRepo
}

func (r *PostRepository) GetPost(ctx context.Context, id entity.UUID) (*entity.Post, error) {
	p := &entity.Post{
		ID: id,
	}
	if err := r.db.WithContext(ctx).Model(&p).WherePK().Select(); err != nil {
		return nil, err
	}
	return p, nil
}

type IPostRepository interface {
	GetPost(ctx context.Context, id entity.UUID) (*entity.Post, error)
}
