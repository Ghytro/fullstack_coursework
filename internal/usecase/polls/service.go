package polls

import (
	"context"
	"fmt"
	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/repository"
	"github.com/Ghytro/galleryapp/internal/validation"
	"github.com/Ghytro/galleryapp/internal/view/polls"
	"log"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/samber/lo"
)

type trendingCache struct {
	Lock      sync.Mutex
	ExpiresAt time.Time
	Polls     []*polls.TrendingPoll
	IsUpdated int32
}

func (c *trendingCache) GetPage(page, pageSize int) []*polls.TrendingPoll {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	startIdx, endIdx := (page-1)*pageSize, page*pageSize
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(c.Polls) {
		endIdx = len(c.Polls)
	}
	result := make([]*polls.TrendingPoll, endIdx-startIdx)
	copy(result, c.Polls[startIdx:endIdx])
	return result
}

// Update не копирует список новых голосов а просто присваивает ссылку
func (c *trendingCache) Update(p []*polls.TrendingPoll) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.Polls = p
	c.ExpiresAt = time.Now().Add(time.Hour)
}

type Service struct {
	repo               Repository
	trendingPollsCache trendingCache
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		trendingPollsCache: trendingCache{
			ExpiresAt: time.Now().Add(-time.Second), // чтобы сразу протух
			Polls:     make([]*polls.TrendingPoll, 0),
		},
	}
}

func (s *Service) CreatePoll(ctx context.Context, creatorID entity.PK, model *polls.NewPollRequest) (*entity.Poll, error) {
	if err := validation.ValidateCreatedPoll(model); err != nil {
		return nil, err
	}
	p := &entity.Poll{
		CreatorID:      creatorID,
		Topic:          model.Topic,
		IsAnonymous:    model.IsAnonymous == "on",
		MultipleChoice: model.MultipleChoice == "on",
		RevoteAbility:  model.CantRevote != "on",
	}
	for i, o := range model.Options {
		p.Options = append(p.Options, &entity.PollOption{
			Index:  i + 1,
			Option: o,
		})
	}
	err := s.repo.CreatePoll(ctx, p)
	return p, err
}

func (s *Service) GetPollWithVotesAmount(ctx context.Context, id entity.PK, userID entity.PK) (*entity.Poll, []*entity.Vote, error) {
	var (
		p         *entity.Poll
		userVotes []*entity.Vote
	)
	err := s.repo.RunInTransaction(ctx, func(tx *database.TX) error {
		repo := s.repo.WithTX(tx)
		var err error
		p, err = repo.GetPoll(ctx, id)
		if err != nil {
			return err
		}
		opts, err := repo.GetVotesAmount(ctx, id)
		if err != nil {
			return err
		}
		p.Options = opts
		creator, err := repo.GetPollCreator(ctx, id)
		if err != nil {
			return err
		}
		p.Creator = creator
		userVotes, err = repo.GetUserPollVotes(ctx, userID, id)
		if err == pg.ErrNoRows {
			err = nil
		}
		return err
	})
	return p, userVotes, err
}

func (s *Service) Vote(ctx context.Context, userID entity.PK, pollID entity.PK, optIdxs ...int) error {
	if len(optIdxs) == 0 {
		return nil
	}
	return s.repo.RunInTransaction(ctx, func(tx *database.TX) error {
		repo := s.repo.WithTX(tx)
		return repo.Vote(ctx, userID, pollID, optIdxs...)
	})
}

func (s *Service) Unvote(ctx context.Context, userID entity.PK, pollID entity.PK) error {
	return s.repo.Unvote(ctx, userID, pollID)
}

func (s *Service) GetMyPolls(ctx context.Context, userID entity.PK, page, pageSize int) ([]*entity.Poll, error) {
	return s.repo.GetPollListSearch(ctx, &repository.PollSearchFilter{
		CreatorID: userID,
	})
}

func (s *Service) updateCache() {
	currTime := time.Now()
	atomic.StoreInt32(&s.trendingPollsCache.IsUpdated, 1)
	defer func() {
		atomic.StoreInt32(&s.trendingPollsCache.IsUpdated, 0)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var trendingPolls []*polls.TrendingPoll
	if err := s.repo.RunInTransaction(ctx, func(tx *database.TX) error {
		repo := s.repo.WithTX(tx)
		votes, err := repo.GetVoteListSearch(ctx, &repository.VoteSearchFilter{
			CreatedAt: &common.Range[time.Time]{
				From: currTime.Add(-(24 * time.Hour)),
				To:   currTime,
			},
		})
		if err != nil {
			return err
		}
		m := make(map[entity.PK]int) // посчитаем количество самых популярных опросов
		for _, v := range votes {
			m[v.PollID]++
		}
		// отсортируем
		pollsWithVotes := lo.MapToSlice(m, func(key entity.PK, value int) struct {
			PollID entity.PK
			Amount int
		} {
			return struct {
				PollID entity.PK
				Amount int
			}{
				PollID: key,
				Amount: value,
			}
		})
		pollIds := lo.Map(pollsWithVotes, func(el struct {
			PollID entity.PK
			Amount int
		}, _ int) entity.PK {
			return el.PollID
		})
		pollsList, err := repo.GetPollListSearch(ctx, &repository.PollSearchFilter{
			IDs: pollIds,
		})
		if err != nil {
			return err
		}
		trendingPolls = lo.Map(pollsList, func(p *entity.Poll, _ int) *polls.TrendingPoll {
			return &polls.TrendingPoll{
				Poll:       p,
				VoteAmount: m[p.ID],
			}
		})
		sort.Slice(trendingPolls, func(i, j int) bool {
			return trendingPolls[i].VoteAmount > trendingPolls[j].VoteAmount
		})
		return nil
	}); err != nil {
		log.Printf("unable to update poll cache: %v\n", err)
		s.trendingPollsCache.ExpiresAt = currTime.Add(time.Second * 10) // попробуй еще раз через 10 секунд
		return
	}
	s.trendingPollsCache.Update(trendingPolls)
}

func (s *Service) GetTrending(ctx context.Context, pageNumber, pageSize int) ([]*polls.TrendingPoll, error) {
	currTime := time.Now()
	result := s.trendingPollsCache.GetPage(pageNumber, pageSize)
	fmt.Println(result)
	if currTime.After(s.trendingPollsCache.ExpiresAt) && atomic.LoadInt32(&s.trendingPollsCache.IsUpdated) == 0 {
		go s.updateCache()
	}
	return result, nil
}
