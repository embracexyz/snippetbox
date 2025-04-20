package main

import (
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
