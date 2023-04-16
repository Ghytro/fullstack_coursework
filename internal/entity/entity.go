package entity

import (
	"errors"
	"time"

	"github.com/Ghytro/galleryapp/internal/database/objectstore"
	"github.com/Ghytro/galleryapp/internal/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PK uint

type ImageID objectstore.FileID

type User struct {
	tableName struct{} `pg:"users"`

	ID            PK       `pg:"id,pk" json:"id"`
	Username      string   `pg:"username,notnull" json:"username"`
	Password      string   `pg:"password,notnull" json:"password"`
	Bio           *string  `pg:"bio" json:"bio"`
	AvatarImageID *ImageID `pg:"avatar_image_id,type:char(24)" json:"avatar_image_id"`
}

func (user *User) Validate() error {
	if err := validation.ValidateUserName(user.Username); err != nil {
		return err
	}

	if err := validation.ValidateUserPassword(user.Password); err != nil {
		return err
	}
	return nil
}

type UUID uuid.UUID

type Post struct {
	tableName struct{} `pg:"posts"`

	ID UUID `pg:"id,pk"`

	UserID PK    `pg:"user_id" json:"-"`
	User   *User `pg:"rel:has-one" json:"-"`

	Caption *string `pg:"caption" json:"caption"`
	ImageID ImageID `pg:"image_id,type:char(24)" json:"image_id"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`

	Comments    []*Comment `pg:"rel:has-many"`
	LikesAmount uint       `json:"likes"`
	Likes       []*Like    `pg:"rel:has-many" json:"who_liked"`
}

type Like struct {
	tableName struct{} `pg:"likes"`

	CreatedAt time.Time `pg:"created_at"`

	PostID UUID  `pg:"post_id"`
	Post   *Post `pg:"rel:has-one"`

	UserID PK    `pg:"user_id"`
	User   *User `pg:"rel:has-one"`
}

type Comment struct {
	tableName struct{} `pg:"comments"`

	ID UUID `pg:"id,pk"`

	CreatedAt time.Time `pg:"created_at"`

	PostID UUID  `pg:"post_id"`
	Post   *Post `pg:"rel:has-one"`

	UserID PK    `pg:"user_id"`
	User   *User `pg:"rel:has-one"`

	Content string `pg:"content,use_zero"`
}

type ErrResponse struct {
	StatusCode int
	Err        error
}

func (e *ErrResponse) Error() string {
	return e.Err.Error()
}

func (e *ErrResponse) Unwrap() error {
	return e.Err
}

func ErrRespBadRequest(err error) *ErrResponse {
	return &ErrResponse{
		StatusCode: fiber.StatusBadRequest,
		Err:        err,
	}
}

func ErrRespNotFound(err error) *ErrResponse {
	return &ErrResponse{
		StatusCode: fiber.StatusNotFound,
		Err:        err,
	}
}

func ErrRespInternalServerError(err error) *ErrResponse {
	return &ErrResponse{
		StatusCode: fiber.StatusInternalServerError,
		Err:        err,
	}
}

func ErrRespUnauthorized(err error) *ErrResponse {
	return &ErrResponse{
		StatusCode: fiber.StatusUnauthorized,
		Err:        err,
	}
}

func ErrRespIncorrectForm() *ErrResponse {
	return ErrRespBadRequest(errors.New("не удалось раскодировать форму"))
}
