package polls

import (
	"errors"
	"fmt"
	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/view"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type View struct {
	service UseCase
}

func NewView(s UseCase) *View {
	return &View{
		service: s,
	}
}

func (v *View) Routers(router fiber.Router, authHandler fiber.Handler, middlewares ...fiber.Handler) {
	r := fiber.New()
	for _, m := range middlewares {
		r.Use(m)
	}
	r.Use(authHandler)
	r.Get("/new/", v.newPollPage)
	r.Get("/trending/:page", v.getTrending)
	r.Get("/my/:page", v.getMyPolls)
	r.Get("/:id", v.getPoll)
	r.Post("/:id/vote", v.vote)
	r.Post("/", v.postNewPoll)
	r.Post("/:id/unvote", v.unvote)
	router.Mount("/polls", r)
}

func (v *View) newPollPage(c *fiber.Ctx) error {
	tpl := templates.MustGet("polls/new.html")
	return view.SendTemplate(c, tpl, nil)
}

func (v *View) postNewPoll(c *fiber.Ctx) error {
	var model NewPollRequest
	if err := c.BodyParser(&model); err != nil {
		return entity.ErrRespBadRequest(err)
	}
	model.IsAnonymous = c.FormValue("is_anonymous")
	model.MultipleChoice = c.FormValue("multiple_choice")
	model.CantRevote = c.FormValue("cant_revote")
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("авторизуйтесь заново"))
	}
	poll, err := v.service.CreatePoll(c.Context(), user.ID, &model)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	return c.Redirect(fmt.Sprintf("/polls/%d", poll.ID), fiber.StatusSeeOther)
}

func (v *View) getPoll(c *fiber.Ctx) error {
	pollID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("авторизуйтесь заново"))
	}
	poll, currentUserVotes, err := v.service.GetPollWithVotesAmount(c.Context(), entity.PK(pollID), user.ID)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	viewData := &GetPollViewData{
		PollID:           fmt.Sprint(poll.ID),
		Topic:            poll.Topic,
		UserID:           fmt.Sprint(poll.Creator.ID),
		Username:         poll.Creator.Username,
		IsAnonymous:      poll.IsAnonymous,
		MultipleChoice:   poll.MultipleChoice,
		RevoteAbility:    poll.RevoteAbility,
		CurrentUserVoted: len(currentUserVotes) != 0,
	}
	for _, o := range poll.Options {
		viewData.Options = append(viewData.Options, Option{
			Option:      o.Option,
			VotesNumber: fmt.Sprint(o.VotesAmount),
		})
	}

	if viewData.CurrentUserVoted {
		viewData.CurrentUserVotes = make([]bool, len(poll.Options))
		for _, v := range currentUserVotes {
			_, optIdx, _ := lo.FindIndexOf(poll.Options, func(o *entity.PollOption) bool {
				return o.ID == v.OptionID
			})
			viewData.CurrentUserVotes[optIdx] = true
		}
	}

	tpl := templates.MustGet("polls/get.html")
	return view.SendTemplate(c, tpl, viewData)
}

func (v *View) vote(c *fiber.Ctx) error {
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("авторизуйтесь заново"))
	}
	pollID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	var model VoteRequest
	if err := c.BodyParser(&model); err != nil {
		return entity.ErrRespBadRequest(err)
	}
	idxs := lo.Map(model.VotesIdxs, func(strIdx string, _ int) int {
		i, err := strconv.Atoi(strIdx)
		if err != nil {
			i = -1
		}
		return i
	})
	if lo.Contains(idxs, -1) {
		return entity.ErrRespBadRequest(errors.New("значение выбранной опции может быть только числовым"))
	}
	if err := v.service.Vote(c.Context(), user.ID, entity.PK(pollID), idxs...); err != nil {
		return entity.ErrRespBadRequest(err)
	}
	return c.Redirect(fmt.Sprintf("/polls/%d", pollID), fiber.StatusSeeOther)
}

func (v *View) unvote(c *fiber.Ctx) error {
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("авторизуйтесь заново"))
	}
	pollID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	if err := v.service.Unvote(c.Context(), user.ID, entity.PK(pollID)); err != nil {
		return entity.ErrRespBadRequest(err)
	}
	return c.Redirect(fmt.Sprintf("/polls/%d", pollID), fiber.StatusSeeOther)
}

func (v *View) getMyPolls(c *fiber.Ctx) error {
	user, ok := c.Locals("user_entity").(*entity.User)
	if !ok {
		return entity.ErrRespUnauthorized(errors.New("авторизуйтесь заново"))
	}

	page, err := strconv.Atoi(c.Params("page"))
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	if page <= 0 {
		return entity.ErrRespBadRequest(errors.New("некорректное значение номера страницы"))
	}

	pageSize, err := strconv.Atoi(c.Query("psize"))
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	if pageSize < 1 {
		return entity.ErrRespBadRequest(errors.New("некорректное значение размера страницы"))
	}

	polls, err := v.service.GetMyPolls(c.Context(), user.ID, page, pageSize)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}

	nextPolls, err := v.service.GetMyPolls(c.Context(), user.ID, page+1, pageSize)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}

	prevPage := "1"
	if page > 1 {
		prevPage = fmt.Sprint(page - 1)
	}
	nextPage := fmt.Sprint(page)
	if len(nextPolls) > 0 {
		nextPage = fmt.Sprint(page + 1)
	}

	tpl := templates.MustGet("polls/my.html")
	viewData := &GetMyPollsViewData{
		PageNumber:     fmt.Sprint(page),
		PageSize:       fmt.Sprint(pageSize),
		PrevPageNumber: prevPage,
		NextPageNumber: nextPage,
		Polls: lo.Map(polls, func(p *entity.Poll, _ int) Poll {
			return Poll{
				ID:             fmt.Sprint(p.ID),
				Title:          p.Topic,
				CreatedAt:      p.CreatedAt.Format("02.01.2006 15:04:05"),
				IsAnonymous:    p.IsAnonymous,
				RevoteAbility:  p.RevoteAbility,
				MultipleChoice: p.MultipleChoice,
				Options: lo.Map(p.Options, func(opt *entity.PollOption, _ int) string {
					return opt.Option
				}),
			}
		}),
	}
	return view.SendTemplate(c, tpl, viewData)
}

func (v *View) getTrending(c *fiber.Ctx) error {
	trending, err := v.service.GetTrending(c.Context(), 1, 10)
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}

	tpl := templates.MustGet("polls/trending.html")

	viewData := &TrendingPollsViewData{
		Polls: lo.Map(trending, func(p *TrendingPoll, _ int) *TrendingPollView {
			return &TrendingPollView{
				ID:         fmt.Sprint(p.ID),
				VoteAmount: fmt.Sprint(p.VoteAmount),
				Title:      p.Topic,
				Options: lo.Map(p.Options, func(opt *entity.PollOption, _ int) string {
					return opt.Option
				}),
				IsAnonymous:    p.IsAnonymous,
				RevoteAbility:  p.RevoteAbility,
				MultipleChoice: p.MultipleChoice,
			}
		}),
	}
	return view.SendTemplate(c, tpl, viewData)
}
