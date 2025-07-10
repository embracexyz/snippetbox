package main

import (
	"net/http"
	"testing"

	"github.com/embracexyz/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	application := newTestApplication(t)
	ts := newTestServer(t, application.getRoutes())
	defer ts.Close()

	status, _, body := ts.get(t, "/ping")

	assert.Assert(t, status, http.StatusOK)
	assert.Assert(t, body, "OK")
}
