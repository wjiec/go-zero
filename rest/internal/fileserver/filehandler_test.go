package fileserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		dir             string
		requestPath     string
		expectedStatus  int
		expectedContent string
	}{
		{
			name:            "Serve static file",
			path:            "/static/",
			dir:             "./testdata",
			requestPath:     "/static/example.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: "1",
		},
		{
			name:           "Pass through non-matching path",
			path:           "/static/",
			dir:            "./testdata",
			requestPath:    "/other/path",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "Directory with trailing slash",
			path:            "/assets",
			dir:             "testdata",
			requestPath:     "/assets/sample.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := Middleware(tt.path, tt.dir)
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})

			handlerToTest := middleware(nextHandler)

			req := httptest.NewRequest("GET", tt.requestPath, nil)
			rr := httptest.NewRecorder()

			handlerToTest.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if len(tt.expectedContent) > 0 {
				assert.Equal(t, tt.expectedContent, rr.Body.String())
			}
		})
	}
}

func TestEnsureTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"path", "path/"},
		{"path/", "path/"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ensureTrailingSlash(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnsureNoTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"path", "path"},
		{"path/", "path"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ensureNoTrailingSlash(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
