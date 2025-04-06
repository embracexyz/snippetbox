package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// 实现/精准匹配，而非通配
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return // 这里不return，会继续执行后续code
	}
	w.Write([]byte("hello from SnippetBox!"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display a snippet"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create a snippet"))
}

func main() {
	// 一个web server所需部分拆解：
	// 1. http server接受http请求，回复http响应
	// 2. 接受到请求后，需要根据urlPath，路由到对应的处理函数，这里是serverMux（管理一组url和响应函数的映射）
	// 3. 若干响应函数，定义了如何接受处理请求、并做出如何的响应

	mux := http.NewServeMux()

	// 在serverMux中； 以/结尾的为通配url，前缀匹配；否则则是表示精准匹配url
	// fixed path & subtree path
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}

}
