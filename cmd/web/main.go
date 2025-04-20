package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr string
}

// 依赖注入;通过把handlerFunc变成appliction的方法，从而使得各类业务函数能access到appliction的属性，即infoLog，实现依赖注入
// 但！仅限于同package，如果是handler分布在不同packege，只能通过closure方式实现，外部package的handlerFunc接受applaction并返回一个http.HandlerFunc类型，通过closure访问appliciton
type appliction struct {
	infoLog *log.Logger
	errLog  *log.Logger
}

func main() {

	var config config
	flag.StringVar(&config.addr, "addr", ":4000", "addr of server")
	flag.Parse()

	// leveled logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Lshortfile)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Llongfile)

	app := appliction{
		infoLog: infoLog,
		errLog:  errLog,
	}

	// 使用自定义http.Server，而非默认的
	server := &http.Server{
		Addr:     config.addr,
		ErrorLog: errLog,
		Handler:  app.getRoutes(),
	}

	infoLog.Printf("Listening on %s", config.addr)
	err := server.ListenAndServe()
	errLog.Fatal(err)
}
