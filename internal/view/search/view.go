package search

import (
	"fmt"
	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/entity"
	"github.com/Ghytro/galleryapp/internal/usecase/search"
	"github.com/Ghytro/galleryapp/internal/view"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type View struct {
	service UseCase
}

func NewView(service UseCase) *View {
	return &View{
		service: service,
	}
}

func (v *View) Routers(router fiber.Router, authHandler fiber.Handler, middleware ...fiber.Handler) {
	r := fiber.New()
	r.Get("/", v.generalSearch)
	router.Mount("/search", r)
}

func (v View) generalSearch(c *fiber.Ctx) error {
	var model searchRequest
	if err := c.BodyParser(&model); err != nil {
		return entity.ErrRespBadRequest(err)
	}
	pageData := common.PageData{
		Page:     1,
		PageSize: 5,
	}
	result, err := v.service.Search(c.Context(), model.Query, &search.PageData{
		UserPage: pageData,
		PollPage: pageData,
	})
	if err != nil {
		return entity.ErrRespBadRequest(err)
	}
	tpl := templates.MustGet("search/general.html")

	viewData := &GeneralSearchViewData{
		Users: lo.Map(result.Users, func(u *entity.User, _ int) User {
			userCountry := "N/A"
			if u.Country != nil {
				userCountry = *u.Country
			}
			return User{
				ID:       fmt.Sprint(u.ID),
				Username: u.Username,
				Country:  userCountry,
			}
		}),
		Polls: lo.Map(result.Polls, func(p *entity.Poll, _ int) Poll {
			return Poll{
				ID:    fmt.Sprint(p.ID),
				Topic: p.Topic,
			}
		}),
	}
	return view.SendTemplate(c, tpl, viewData)
}
