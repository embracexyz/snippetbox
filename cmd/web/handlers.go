package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/embracexyz/snippetbox/internal/models"
)

func (app *appliction) download(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/static/file.txt")
}

func (app *appliction) home(w http.ResponseWriter, r *http.Request) {
	// 实现/精准匹配，而非通配
	if r.URL.Path != "/" {
		app.notFound(w)
		return // 这里不return，会继续执行后续code
	}

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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
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

func (app *appliction) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
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
