package polls

import (
	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/view"
	"html/template"
)

var templates *common.SyncMap[string, *template.Template]

func init() {
	var err error
	templates, err = view.GenTemplatesMap(
		"polls/new.html",
		"polls/get.html",
		"polls/my.html",
		"polls/trending.html",
	)
	common.LogFatalErr(err)
}
