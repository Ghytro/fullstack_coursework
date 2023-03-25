package entity

import (
	"errors"
	"time"

	"github.com/go-pg/pg/v10/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PK uint

type baseEntity struct {
	CreatedAt time.Time      `pg:"created_at"`
	DeletedAt types.NullTime `pg:"deleted_at"`
}

type pkID struct {
	ID PK `pg:"id,pk"`
}

type User struct {
	tableName struct{} `pg:"users"`

	pkID
	baseEntity

	Username  string  `pg:"username,notnull" form:"username"`
	FirstName *string `pg:"first_name" form:"first_name"`
	LastName  *string `pg:"last_name" form:"last_name"`

	Password  string  `pg:"password,notnull" form:"password"`
	Bio       *string `pg:"bio" form:"bio"`
	AvatarUrl *string `pg:"avatar_url" form:"avatar_url"`
	Country   *string `pg:"country" form:"country"`

	Polls []*Poll `pg:"rel:has-many"`
	Votes []*Vote `pg:"rel:has-many"`
}

type Poll struct {
	tableName struct{} `pg:"polls"`

	pkID
	baseEntity

	CreatorID PK    `pg:"creator_id"`
	Creator   *User `pg:"rel:has-one"`

	Topic          string        `pg:"topic,notnull"`
	IsAnonymous    bool          `pg:"is_anonymous,notnull"`
	MultipleChoice bool          `pg:"multiple_choice,notnull"`
	RevoteAbility  bool          `pg:"revote_ability,notnull"`
	Options        []*PollOption `pg:"rel:has-many"`
}

type PollOption struct {
	tableName struct{} `pg:"poll_options"`

	pkID

	PollID PK    `pg:"poll_id"`
	Poll   *Poll `pg:"rel:has-one"`

	Votes []*Vote `pg:"rel:has-many"`

	VotesAmount int `pg:"-"`

	Index     int            `pg:"index,notnull"`
	Option    string         `pg:"option,notnull"`
	UpdatedAt types.NullTime `pg:"updated_at"`
}

type Vote struct {
	tableName struct{} `pg:"votes"`

	baseEntity
	ID uuid.UUID `pg:"id"`

	UserID PK    `pg:"user_id,notnull"`
	User   *User `pg:"rel:has-one"`

	OptionID PK          `pg:"option_id,notnull"`
	Option   *PollOption `pg:"rel:has-one"`

	PollID PK    `pg:"poll_id,notnull"`
	Poll   *Poll `pg:"rel:has-one"`
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
