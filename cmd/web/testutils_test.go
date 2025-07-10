package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/embracexyz/snippetbox/internal/models/mock"
	"github.com/go-playground/form/v4"
)

func newTestApplication(t *testing.T) *application {

	formDecoder := form.NewDecoder()

	scs := scs.New()
	scs.Lifetime = 12 * time.Hour
	scs.Cookie.Secure = true

	templateCache, err := NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		errLog:         log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		users:          &mock.MockUserModel{},
		snippets:       &mock.MockSnippetModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: scs,
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// client后续收到的cookie都存在jar里，并会在后续请求时带上
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar

	// 收到redirect请求时，立刻抛err而不是进行跳转
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

type testServer struct {
	*httptest.Server
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}
