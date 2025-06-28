package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/embracexyz/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *appliction) download(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/static/file.txt")
}

func (app *appliction) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := NewTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *appliction) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := NewTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)

}

func (app *appliction) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// r.Form : 包含任意请求方法的请求体参数、和query 参数
	// r.PostForm: 只包括post/patch/put方法的的请求体参数
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validation before insert

	form := Form{
		Title:       title,
		Content:     content,
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 10 {
		form.FieldErrors["title"] = "This field is too long (maximum is 10 characters)"
	}
	if strings.TrimSpace(content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(content) > 5 {
		form.FieldErrors["content"] = "This field is too long (maximum is 5 characters)"
	}
	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *appliction) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData(r)
	data.Form = Form{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

// 捕捉panic的中间件只能补充同一个goroutine的panice，如果某handler启动了一个新的goroutine，且在其中panic了，那么中间件就无法handle，因此：需要在call 另一个goroutine之前，自己注册一个recover的defer
func (app *appliction) spinNewGoroutineHandler(w http.ResponseWriter, r *http.Request) {

	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.infoLog.Println("recover in sub goroutine")
			}
		}()

		// do some background stuff which may be panic
	}()

	w.Write([]byte("ok~"))
}
