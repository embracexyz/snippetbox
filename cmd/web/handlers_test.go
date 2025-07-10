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

func TestSnippetView(t *testing.T) {
	application := newTestApplication(t)
	ts := newTestServer(t, application.getRoutes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Assert(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}

}
