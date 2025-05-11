package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/embracexyz/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func NewTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// 提供一次性加载template到内存的方法，减少多次磁盘io
func NewTemplateCache() (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page) // 返回文件名, 以文件名做key

		// ts, err := template.ParseFiles("./ui/html/pages/base.tmpl")
		ts, err := template.New(fileName).Funcs(functions).ParseFiles("./ui/html/pages/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		templates[fileName] = ts
	}

	return templates, nil
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
