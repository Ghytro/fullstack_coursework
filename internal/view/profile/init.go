package profile

import (
	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/view"
	"html/template"
)

var templates *common.SyncMap[string, *template.Template]

func init() {
	var err error
	templates, err = view.GenTemplatesMap(
		"profile/any.html",
		"profile/my.html",
	)
	common.LogFatalErr(err)
}
