package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"

	"github.com/embracexyz/snippetbox/internal/models"
)

type config struct {
	addr string
	dsn  string
}

// 依赖注入;通过把handlerFunc变成application的方法，从而使得各类业务函数能access到application的属性，即infoLog，实现依赖注入
// 但！仅限于同package，如果是handler分布在不同packege，只能通过closure方式实现，外部package的handlerFunc接受applaction并返回一个http.HandlerFunc类型，通过closure访问appliciton
type application struct {
	infoLog        *log.Logger
	errLog         *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	var config config
	flag.StringVar(&config.addr, "addr", ":4000", "addr of server")
	flag.StringVar(&config.dsn, "dsn", "web:yourpassword@(127.0.0.1:13306)/snippetbox?parseTime=true", "mysql datasource name")
	flag.Parse()

	// leveled logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Lshortfile)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Llongfile)

	db, err := openDB(config.dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	// init template cache
	templateCache, err := NewTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	// add formDecoder
	formDecoder := form.NewDecoder()

	// add sessionManager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true // for https connection

	// init application

	app := application{
		infoLog:        infoLog,
		errLog:         errLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// 使用自定义http.Server，而非默认的
	server := &http.Server{
		Addr:     config.addr,
		ErrorLog: errLog,
		Handler:  app.getRoutes(),
		// go http server默认启动了keep alive，即多个http请求可以复用一个tcp连接，减少tcp建立释放的开销，但是当tcp idle一定时间server端要主动断开，避免slow client的攻击笑话服务端资源
		IdleTimeout: time.Minute,
		// 如果超过5s，请求头或者请求体还没读完，就关闭tcp连接
		ReadTimeout: 5 * time.Second,
		// http协议：从开始读request header时开始计时，https协议，从requst被接受开始计时
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Listening on %s", config.addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errLog.Fatal(err)
}
