package search

import (
	"context"
	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/usecase/search"
)

type UseCase interface {
	Search(ctx context.Context, query string, page *search.PageData) (*search.SearchResult, error)
	SearchUser(ctx context.Context, searchParams *search.UserSearchParams, page *common.PageData) ([]*entity.User, error)
	SearchPoll(ctx context.Context, searchParams *search.PollSearchParams, page *common.PageData) ([]*entity.Poll, error)
}
