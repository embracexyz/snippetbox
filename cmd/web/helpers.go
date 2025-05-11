package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *appliction) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *appliction) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *appliction) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// 抽离render template 的通用逻辑
func (app *appliction) render(w http.ResponseWriter, status int, templateName string, data *templateData) {
	ts, ok := app.templateCache[templateName]
	if !ok {
		err := fmt.Errorf("template %s not found!", templateName)
		app.infoLog.Fatalf("%+v", app.templateCache)
		app.serverError(w, err)
		return
	}

	// catch runtime errors
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
