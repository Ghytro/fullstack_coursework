package repository

import (
	"context"
	"errors"

	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/entity"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/samber/lo"
)

type UserRepository struct {
	db DBI
}

func NewUserRepo(db DBI) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (m *UserRepository) RunInTransaction(ctx context.Context, f func(tx *database.TX) error) error {
	return m.db.RunInTransaction(ctx, f)
}

func (m *UserRepository) WithTX(tx *database.TX) *UserRepository {
	return NewUserRepo(tx)
}

func (m *UserRepository) CreateUser(ctx context.Context, user *entity.User) (entity.PK, error) {
	var u entity.User
	err := m.db.RunInTransaction(ctx, func(tx *database.TX) error {
		if err := tx.ModelContext(ctx, &u).Where("username = ?", user.Username).Select(); err != nil {
			if err == pg.ErrNoRows {
				_, _err := tx.ModelContext(ctx, &u).
					Value("username", "?", user.Username).
					Value("password", "crypt(?, gen_salt('bf'))", user.Password).
					Value("first_name", "?", user.FirstName).
					Value("last_name", "?", user.LastName).
					Value("bio", "?", user.Bio).
					Value("avatar_url", "?", user.AvatarUrl).
					Value("country", "?", user.Country).
					Returning("*").Insert()
				return _err
			}
			return err
		}
		if u.DeletedAt.IsZero() {
			return errors.New("пользователь с таким именем уже существует")
		}
		_, err := tx.ModelContext(ctx, &u).Where("id = ?", u.ID).Set(
			`username = ?,
			first_name = ?,
			last_name = ?,
			password = crypt(?, password),
			bio = ?,
			avatar_url = ?,
			country = ?,
			created_at = NOW(),
			deleted_at = NULL`,
			user.Username,
			user.FirstName,
			user.LastName,
			user.Password,
			user.Bio,
			user.AvatarUrl,
			user.Country,
		).Update()
		return err
	})
	return u.ID, err
}

func (m *UserRepository) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	var u entity.User
	if err := m.db.ModelContext(ctx, &u).Where("id = ? AND deleted_at IS NULL", userID).Select(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *UserRepository) GetUserWithPolls(ctx context.Context, userID entity.PK, limit int) (*entity.User, error) {
	var (
		u       entity.User
		options []*entity.PollOption
	)
	err := m.db.RunInTransaction(ctx, func(tx *database.TX) error {
		if err := tx.ModelContext(ctx, &u).Where("id = ? AND deleted_at IS NULL", userID).Select(); err != nil {
			return err
		}
		q := tx.ModelContext(ctx, &u.Polls).Where("creator_id = ?", userID).Order("created_at DESC")
		if limit != -1 {
			q = q.Limit(limit)
		}
		if err := q.Select(); err != nil {
			return err
		}
		if len(u.Polls) == 0 {
			return nil
		}
		pollsIDs := lo.Map(u.Polls, func(item *entity.Poll, index int) entity.PK {
			return item.ID
		})
		return tx.ModelContext(ctx, &options).Where("poll_id IN (?)", pg.In(pollsIDs)).Order("index ASC").Select()
	})
	if err != nil {
		return nil, err
	}
	pollsMapping := make(map[entity.PK]*entity.Poll)
	for _, p := range u.Polls {
		pollsMapping[p.ID] = p
	}
	for _, o := range options {
		p := pollsMapping[o.PollID]
		p.Options = append(p.Options, o)
	}
	return &u, nil
}

func (m *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	_, err := m.db.ModelContext(ctx, user).Where("id = ? AND deleted_at IS NULL", user.ID).Update()
	return err
}

func (m *UserRepository) DeleteUser(ctx context.Context, userID entity.PK) error {
	_, err := m.db.ModelContext(ctx, (*entity.User)(nil)).Set("deleted_at = NOW()").Where("id = ?", userID).Update()
	return err
}

func (m *UserRepository) Auth(ctx context.Context, username string, password string) (entity.PK, error) {
	var u entity.User
	if err := m.db.ModelContext(ctx, &u).Where("username = ? AND password = crypt(?, password)", username, password).Select(); err != nil {
		return 0, err
	}
	return u.ID, nil
}

type UserSearchFilter struct {
	IDs      []entity.PK
	UserName *common.StringDataFilter
	PageData *common.PageData
}

func (f UserSearchFilter) Apply(q *orm.Query) {
	if f.UserName != nil {
		f.UserName.Apply(q, "username")
	}
	if len(f.IDs) > 0 {
		q.Where("id in (?)", pg.In(f.IDs))
	}
	if f.PageData != nil {
		q.Limit(f.PageData.PageSize).Offset((f.PageData.Page - 1) * f.PageData.PageSize)
	}
}

func (m *UserRepository) GetUserListSearch(ctx context.Context, filter *UserSearchFilter) ([]*entity.User, error) {
	var u []*entity.User
	err := m.db.RunInTransaction(ctx, func(tx *database.TX) error {
		q := tx.ModelContext(ctx, &u)
		filter.Apply(q)
		return q.Select()
	})
	return u, err
}

type IUserRepository interface {
}
