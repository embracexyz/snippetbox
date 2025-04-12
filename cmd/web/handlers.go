package main

import (
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// 实现/精准匹配，而非通配
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return // 这里不return，会继续执行后续code
	}

	// 文件路径为绝对路径，或者相对路径（当前执行go run 的路径）
	ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	if err != nil {
		http.Error(w, "Internal Server Error!", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error!", 500)
		return
	}

}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display a snippet"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create a snippet"))
}
