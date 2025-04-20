package main

import (
	"html/template"
	"net/http"
	"strconv"
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

	files := []string{
		"./ui/html/pages/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// 文件路径为绝对路径，或者相对路径（当前执行go run 的路径）
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *appliction) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	w.Write([]byte("display a snippet"))
}

func (app *appliction) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("create a snippet"))
}
