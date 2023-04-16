package repository

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/database"
	"github.com/Ghytro/galleryapp/internal/entity"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type UserRepository struct {
	*MonoRepo
}

func (m *UserRepository) CreateUser(ctx context.Context, user *entity.User) (entity.PK, error) {
	var u entity.User
	if _, err := m.db.WithContext(ctx).Model(&u).
		Value("username", "?", user.Username).
		Value("password", "crypt(?, gen_salt('bf'))", user.Password).
		Value("bio", "?", user.Bio).
		Value("avatar_image_id", "?", user.AvatarImageID).
		Returning("*").Insert(); err != nil {
		return 0, err
	}
	return u.ID, nil
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
		bio = ?,
		avatar_image_id = ?`,

		user.Username,
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

type IUserRepository interface {
	Auth(ctx context.Context, username string, password string) (entity.PK, error)
	CreateUser(ctx context.Context, user *entity.User) (entity.PK, error)
	DeleteUser(ctx context.Context, userID entity.PK) error
	GetUser(ctx context.Context, userID entity.PK) (*entity.User, error)
	GetUserListSearch(ctx context.Context, filter *UserSearchFilter) ([]*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
}
