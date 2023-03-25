package view

import (
	"github.com/Ghytro/galleryapp/internal/entity"
	"html/template"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SendTemplate(c *fiber.Ctx, tpl *template.Template, data interface{}) error {
	var htmlBuf strings.Builder
	if err := tpl.Execute(&htmlBuf, data); err != nil {
		return &entity.ErrResponse{
			StatusCode: fiber.StatusInternalServerError,
			Err:        err,
		}
	}
	return c.SendString(htmlBuf.String())
}
