package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *appliction) getRoutes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// // 原本应该是Handle方法，注册一个http.Handler接口类型（实现了ServeHttp方法）；但是HandleFunc却能直接注册一个函数；
	// // 背后，是因为HandleFunc是包装，本质是调用HandleFunc(f)；对f做了强制类型转换；而该类型实现了ServeHttp方法，通过直接call f函数实现，而该函数签名整好是满足ServeHttp签名
	// // 是一种adapter适配器模式
	// mux.Handle("/", http.HandlerFunc(app.home))
	// // mux.HandleFunc("/", home)

	// mux.HandleFunc("/snippet/view", app.snippetView)
	// mux.HandleFunc("/snippet/create", app.snippetCreate)

	// // server single file
	// mux.HandleFunc("/download", app.download)

	router.HandlerFunc(http.MethodGet, "/", app.home) // excatly match "/"
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// add a middleware chain containing our 'stanard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, app.secureHeaders)

	// 这里mux本身也是满足了ServeHttp接口，这里只是入口路由功能，通过注册的url，分配到不同Handler接口去处理，最后汇总结果并返回
	// 每个handler是并发处理的，单独的goroutine中，因此要注意race condition
	return standard.Then(router)
}
