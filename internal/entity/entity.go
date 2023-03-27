package entity

import (
	"errors"
	"time"

	"github.com/Ghytro/galleryapp/internal/database/objectstore"
	"github.com/go-pg/pg/v10/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PK uint

type baseEntity struct {
	CreatedAt time.Time      `pg:"created_at"`
	UpdatedAt types.NullTime `pg:"updated_at"`
	DeletedAt types.NullTime `pg:"deleted_at"`
}

type ImageID objectstore.FileID

type User struct {
	tableName struct{} `pg:"users"`

	baseEntity

	ID        PK      `pg:"id,pk" json:"id"`
	Username  string  `pg:"username,notnull" json:"username"`
	FirstName *string `pg:"first_name" json:"first_name"`
	LastName  *string `pg:"last_name" json:"last_name"`

	Password      string   `pg:"password,notnull" json:"password"`
	Bio           *string  `pg:"bio" json:"bio"`
	AvatarImageID *ImageID `pg:"avatar_image_id,type:char(24)" json:"avatar_url"`

	Posts         []*Post         `pg:"rel:has-many"`
	Subscribers   []*Subscription `pg:"rel:has-many"`
	Subscriptions []*Subscription `pg:"rel:has-many"`
}

type UUID uuid.UUID

type Post struct {
	tableName struct{} `pg:"posts"`

	baseEntity

	ID UUID `pg:"id,pk"`

	UserID PK    `pg:"user_id"`
	User   *User `pg:"rel:has-one"`

	Caption *string `pg:"caption"`
	ImageID ImageID `pg:"image_id,type:char(24)" json:"image_id"`
}

type Subscription struct {
	tableName struct{} `pg:"subscriptions"`

	SubscriberID PK    `pg:"subscriber_id"`
	Subscriber   *User `pg:"rel:has-one"`

	PublisherID PK    `pg:"publisher_id"`
	Publisher   *User `pg:"rel:has-one"`
}

type Like struct {
	tableName struct{} `pg:"likes"`

	baseEntity

	PostID UUID  `pg:"post_id"`
	Post   *Post `pg:"rel:has-one"`

	UserID PK    `pg:"user_id"`
	User   *User `pg:"rel:has-one"`
}

type Comment struct {
	tableName struct{} `pg:"comments"`

	baseEntity

	ID UUID `pg:"id,pk"`

	PostID UUID  `pg:"post_id"`
	Post   *Post `pg:"rel:has-one"`

	UserID PK    `pg:"user_id"`
	User   *User `pg:"rel:has-one"`

	Content string `pg:"content,use_zero"`
}

type CommentLike struct {
	tableName struct{} `pg:"comment_likes"`

	baseEntity

	CommentID UUID     `pg:"comment_id"`
	Comment   *Comment `pg:"rel:has-one"`

	UserID PK    `pg:"user_id"`
	User   *User `pg:"rel:has-one"`
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
