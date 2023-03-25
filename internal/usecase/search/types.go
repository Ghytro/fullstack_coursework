package search

import (
	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/entity"
)

type UserSearchParams struct {
	Username          *common.StringDataFilter
	RealName, Country *string
}

type PollSearchParams struct {
	Topic, CreatorUsername                     *common.StringDataFilter
	IsAnonymous, RevoteAbility, MultipleChoice *bool
}

type PageData struct {
	UserPage, PollPage common.PageData
}

type SearchResult struct {
	Users []*entity.User
	Polls []*entity.Poll
}
