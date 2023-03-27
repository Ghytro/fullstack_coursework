package repository

import (
	"context"
	"errors"

	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/samber/lo"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type UserRepository struct {
	*MonoRepo
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
					Value("avatar_image_id", "?", user.AvatarImageID).
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
			avatar_image_id = ?,
			created_at = NOW(),
			deleted_at = NULL`,
			user.Username,
			user.FirstName,
			user.LastName,
			user.Password,
			user.Bio,
			user.AvatarImageID,
		).Update()
		return err
	})
	return u.ID, err
}

func (m *UserRepository) GetUser(ctx context.Context, userID entity.PK) (*entity.User, error) {
	var u entity.User
	if err := m.db.WithContext(ctx).Model(&u).Where("id = ? AND deleted_at IS NULL", userID).Select(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	_, err := m.db.WithContext(ctx).Model(user).Where("id = ? AND deleted_at IS NULL", user.ID).Set(
		`username = ?,
		first_name = ?,
		last_name = ?,
		password = crypt(?, password),
		bio = ?,
		avatar_image_id = ?,
		created_at = NOW(),
		deleted_at = NULL`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Bio,
		user.AvatarImageID,
	).Update()
	return err
}

func (m *UserRepository) DeleteUser(ctx context.Context, userID entity.PK) error {
	_, err := m.db.WithContext(ctx).Model((*entity.User)(nil)).Set("deleted_at = NOW()").Where("id = ?", userID).Update()
	return err
}

func (m *UserRepository) Auth(ctx context.Context, username string, password string) (entity.PK, error) {
	var u entity.User
	if err := m.db.WithContext(ctx).Model(ctx, &u).Where("username = ? AND password = crypt(?, password)", username, password).Select(); err != nil {
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

// GetSubscriptionCount собрать полную инфу о подписках по пользователю (сколько подписчиков и подписок)
func (r *UserRepository) GetSubscriptionCount(ctx context.Context, userID entity.PK) (subscribed uint, subscriptions uint, err error) {
	var subscribedI, subscriptionsI int
	err = r.db.RunInTransaction(ctx, func(tx *database.TX) error {
		if subscribedI, err = tx.ModelContext(ctx, (*entity.Subscription)(nil)).Where("publisher_id = ?", userID).Count(); err != nil {
			return err
		}
		subscriptionsI, err = tx.ModelContext(ctx, (*entity.Subscription)(nil)).Where("subscriber_id = ?", userID).Count()
		return err
	})
	return uint(subscribedI), uint(subscriptionsI), err
}

func (r *UserRepository) GetSubscribers(ctx context.Context, userID entity.PK, pageSize, pageNum uint) ([]entity.PK, error) {
	var res []*entity.Subscription
	if err := r.db.WithContext(ctx).Model(&res).Order("publisher_id").Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Select(); err != nil {
		return nil, err
	}
	return lo.Map(res, func(s *entity.Subscription, _ int) entity.PK { return s.SubscriberID }), nil
}

func (r *UserRepository) GetSubscribed(ctx context.Context, userID entity.PK, pageSize, pageNum uint) ([]entity.PK, error) {
	var res []*entity.Subscription
	if err := r.db.WithContext(ctx).Model(&res).Order("subscriber_id").Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Select(); err != nil {
		return nil, err
	}
	return lo.Map(res, func(s *entity.Subscription, _ int) entity.PK { return s.PublisherID }), nil
}

type IUserRepository interface {
	Auth(ctx context.Context, username string, password string) (entity.PK, error)
	CreateUser(ctx context.Context, user *entity.User) (entity.PK, error)
	DeleteUser(ctx context.Context, userID entity.PK) error
	GetUser(ctx context.Context, userID entity.PK) (*entity.User, error)
	GetUserListSearch(ctx context.Context, filter *UserSearchFilter) ([]*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	GetSubscriptionCount(ctx context.Context, userID entity.PK) (subscribed uint, subscriptions uint, err error)
	GetSubscribers(ctx context.Context, userID entity.PK, pageSize, pageNum uint) ([]entity.PK, error)
	GetSubscribed(ctx context.Context, userID entity.PK, pageSize, pageNum uint) ([]entity.PK, error)
}
