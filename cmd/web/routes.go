package main

import (
	"net/http"

	"github.com/embracexyz/snippetbox/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) getRoutes() http.Handler {
	router := httprouter.New()

	// custom error handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	// router.MethodNotAllowed = xxx : in the same way too

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// // 原本应该是Handle方法，注册一个http.Handler接口类型（实现了ServeHttp方法）；但是HandleFunc却能直接注册一个函数；
	// // 背后，是因为HandleFunc是包装，本质是调用HandleFunc(f)；对f做了强制类型转换；而该类型实现了ServeHttp方法，通过直接call f函数实现，而该函数签名整好是满足ServeHttp签名
	// // 是一种adapter适配器模式
	// mux.Handle("/", http.HandlerFunc(app.home))
	// // mux.HandleFunc("/", home)

	// mux.HandleFunc("/snippet/view", app.snippetView)
	// mux.HandleFunc("/snippet/create", app.snippetCreate)

	// // server single file
	// mux.HandleFunc("/download", app.download)

	// 注册session管理的中间件，在它之后的handler，会被自动处理session数据，
	// session中间件做了什么？
	/*
		1. 检查当前request是否有session，有则在对应store里（这里是mysql，也比如说redis）取出session token（一个token标识一个客户端）
			所对应的session data(之前请求所添加的一些上下文信息，比如是否登录了，访问了哪些页面)，当然也会检查其是否过期；只有没过期的session data 会被attach到context中；
			从而被后面的业务handler使用；业务handler处理中更新了session data ，也会被sessionManager保存到store中，再返回给客户端
	*/
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))       // excatly match "/"
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about)) // excatly match "/"
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	// user
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUp))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// ping
	router.Handler(http.MethodGet, "/ping", http.HandlerFunc(ping))

	protected := dynamic.Append(app.requireAuth)
	// add requireAuth middler before this url
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogout))

	// add a middleware chain containing our 'stanard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, app.secureHeaders)

	// 这里mux本身也是满足了ServeHttp接口，这里只是入口路由功能，通过注册的url，分配到不同Handler接口去处理，最后汇总结果并返回
	// 每个handler是并发处理的，单独的goroutine中，因此要注意race condition
	return standard.Then(router)
}
