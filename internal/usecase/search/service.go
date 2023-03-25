package search

import (
	"context"
	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/repository"

	"github.com/go-pg/pg/v10"
)

type Service struct {
	userRepo  UserRepository
	pollsRepo PollsRepository
}

func NewService(userRepo UserRepository, pollsRepo PollsRepository) *Service {
	return &Service{
		userRepo:  userRepo,
		pollsRepo: pollsRepo,
	}
}

func (s Service) Search(ctx context.Context, query string, page *PageData) (*SearchResult, error) {
	var result *SearchResult
	err := s.userRepo.RunInTransaction(ctx, func(tx *database.TX) error {
		userRepo, pollsRepo := s.userRepo.WithTX(tx), s.pollsRepo.WithTX(tx)
		likeFilter := &common.StringDataFilter{
			Likeness: common.StrDataLikenessSubstr,
			Value:    query,
		}
		userResults, err := userRepo.GetUserListSearch(ctx, &repository.UserSearchFilter{
			UserName: likeFilter,
			PageData: &page.UserPage,
		})
		if err != nil && err != pg.ErrNoRows {
			return err
		}
		pollsResults, err := pollsRepo.GetPollListSearch(ctx, &repository.PollSearchFilter{
			Topic:    likeFilter,
			PageData: &page.PollPage,
		})
		if err != nil && err != pg.ErrNoRows {
			return err
		}
		result = &SearchResult{
			Users: userResults,
			Polls: pollsResults,
		}
		return nil
	})
	return result, err
}

func (s Service) SearchPoll(ctx context.Context, searchParams *PollSearchParams, page *common.PageData) ([]*entity.Poll, error) {
	var p []*entity.Poll
	err := s.pollsRepo.RunInTransaction(ctx, func(tx *database.TX) error {
		pollsRepo := s.pollsRepo.WithTX(tx)
		var err error
		p, err = pollsRepo.GetPollListSearch(ctx, &repository.PollSearchFilter{
			Topic:           searchParams.Topic,
			CreatorUserName: searchParams.CreatorUsername,
			IsAnonymous:     searchParams.IsAnonymous,
			MultipleChoice:  searchParams.MultipleChoice,
			RevoteAbility:   searchParams.RevoteAbility,
			PageData:        page,
		})
		return err
	})
	return p, err
}

func (s Service) SearchUser(ctx context.Context, searchParams *UserSearchParams, page *common.PageData) ([]*entity.User, error) {
	var u []*entity.User
	err := s.userRepo.RunInTransaction(ctx, func(tx *database.TX) error {
		userRepo := s.userRepo.WithTX(tx)
		var err error
		u, err = userRepo.GetUserListSearch(ctx, &repository.UserSearchFilter{
			UserName: searchParams.Username,
			PageData: page,
		})
		return err
	})
	return u, err
}
