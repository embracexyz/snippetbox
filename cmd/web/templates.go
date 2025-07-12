package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/embracexyz/snippetbox/internal/models"
	"github.com/embracexyz/snippetbox/internal/validator"
	"github.com/embracexyz/snippetbox/ui"
	"github.com/justinas/nosurf"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	// add a form 存储提交表单的数据和error校验信息
	Form interface{}
	// add Flash
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	About           string
}

type Form struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) NewTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"), // 自动尝试从当前context中取flash信息
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

// 提供一次性加载template到内存的方法，减少多次磁盘io
func NewTemplateCache() (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page) // 返回文件名, 以文件名做key

		patterns := []string{
			"html/pages/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(fileName).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		templates[fileName] = ts
	}

	return templates, nil
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// convert to UTC before formatting
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
