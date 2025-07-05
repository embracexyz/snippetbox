package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/embracexyz/snippetbox/internal/models"
	"github.com/embracexyz/snippetbox/internal/validator"
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
