package view

import (
	"github.com/Ghytro/galleryapp/internal/common"
	"html/template"
	"sync"
)

func GenTemplatesMap(paths ...string) (*common.SyncMap[string, *template.Template], error) {
	m := common.NewSyncMap[string, *template.Template](&sync.Mutex{})
	for _, p := range paths {
		t, err := template.ParseFiles("./web/" + p)
		common.LogFatalErr(err)
		m.Set(p, t)
	}
	return m, nil
}
