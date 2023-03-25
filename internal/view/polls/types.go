package polls

import "github.com/Ghytro/galleryapp/internal/entity"

type NewPollRequest struct {
	Topic          string   `form:"topic"`
	Options        []string `form:"options"`
	IsAnonymous    string   `form:"is_anonymous"`
	MultipleChoice string   `form:"multiple_choice"`
	CantRevote     string   `form:"cant_revote"`
}

type VoteRequest struct {
	VotesIdxs []string `form:"votes"`
}

type GetPollViewData struct {
	PollID,
	Topic,
	UserID,
	Username string
	MultipleChoice,
	RevoteAbility,
	IsAnonymous,
	CurrentUserVoted bool
	Options          []Option
	CurrentUserVotes []bool
}

type Option struct {
	Option, VotesNumber string
}

type GetMyPollsViewData struct {
	PageNumber, PageSize           string
	PrevPageNumber, NextPageNumber string
	Polls                          []Poll
}

type Poll struct {
	ID,
	CreatedAt,
	Title string
	IsAnonymous, RevoteAbility, MultipleChoice bool

	Options []string
}

type TrendingPoll struct {
	*entity.Poll
	VoteAmount int
}

type TrendingPollsViewData struct {
	Polls []*TrendingPollView
}

type TrendingPollView struct {
	ID, Title, CreatedAt, VoteAmount string

	Options []string

	IsAnonymous, RevoteAbility,
	MultipleChoice bool
}
