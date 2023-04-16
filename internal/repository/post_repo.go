package repository

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/util"
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

func (r *PostRepository) GetPostComments(ctx context.Context, id entity.UUID, page util.PaginationData) ([]*entity.Comment, error) {
	return nil, nil // TODO
}

func (r *PostRepository) GetPostsByUser(ctx context.Context, userID entity.PK, page util.PaginationData) (posts []*entity.Post, err error) {
	if err := r.db.WithContext(ctx).
		Model(&posts).
		Where("user_id = ?", userID).
		Limit(int(page.Size)).
		Offset(int(page.Page * page.Size)).
		Order("created_at DESC").
		Select(); err != nil {
		return nil, err
	}
	return
}

type IPostRepository interface {
	GetPost(ctx context.Context, id entity.UUID) (*entity.Post, error)
	GetPostComments(ctx context.Context, id entity.UUID, page util.PaginationData) ([]*entity.Comment, error)
	GetPostsByUser(ctx context.Context, userID entity.PK, page util.PaginationData) ([]*entity.Post, error)
}
