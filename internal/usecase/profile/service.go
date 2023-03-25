package profile

import (
	"context"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/entity"
)

type Service struct {
	repo      ProfileRepo
	pollsRepo PollsRepo
}

func NewService(repo ProfileRepo, pollsRepo PollsRepo) *Service {
	return &Service{
		repo:      repo,
		pollsRepo: pollsRepo,
	}
}

func (s *Service) CreateUser(ctx context.Context, user *entity.User) (entity.PK, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *Service) GetUserWithPolls(ctx context.Context, userID entity.PK, limit int) (*entity.User, error) {
	var u *entity.User
	err := s.repo.RunInTransaction(ctx, func(tx *database.TX) error {
		repo := s.repo.WithTX(tx)
		pollsRepo := s.pollsRepo.WithTX(tx)

		var err error
		u, err = repo.GetUser(ctx, userID)
		if err != nil {
			return err
		}
		polls, err := pollsRepo.GetPollsCreatedBy(ctx, userID, limit, 0)
		if err != nil {
			return err
		}
		u.Polls = polls
		return nil
	})
	return u, err
}

func (s *Service) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, userID entity.PK) error {
	return s.repo.DeleteUser(ctx, userID)
}
