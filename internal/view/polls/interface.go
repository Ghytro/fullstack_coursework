package polls

import (
	"context"
	"github.com/Ghytro/galleryapp/internal/entity"
)

type UseCase interface {
	CreatePoll(ctx context.Context, creatorID entity.PK, model *NewPollRequest) (*entity.Poll, error)
	GetPollWithVotesAmount(ctx context.Context, id entity.PK, userID entity.PK) (*entity.Poll, []*entity.Vote, error)
	Vote(ctx context.Context, userID entity.PK, pollID entity.PK, optIdxs ...int) error
	Unvote(ctx context.Context, userID entity.PK, pollID entity.PK) error
	GetMyPolls(ctx context.Context, userID entity.PK, page, pageSize int) ([]*entity.Poll, error)
	GetTrending(ctx context.Context, pageNumber, pageSize int) ([]*TrendingPoll, error)
}
