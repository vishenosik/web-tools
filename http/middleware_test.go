package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vishenosik/web/versions"
)

func TestApiVersionMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		request        func() *http.Request
		expectError    bool
		expectedStatus int
	}{
		{
			name: "Valid request",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=2.1", nil)
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid version",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=invalid", nil)
				return req
			},
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handlerCalled := false

			middleware := ApiVersionMiddleware(versions.NewDotVersion("2.1"))
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				_, err := TypedApiVersionFromContext[versions.DotVersion](r.Context())
				if err != nil {
					t.Errorf("Failed to get version from context: %v", err)
				}
			})

			middleware(testHandler).ServeHTTP(rr, tt.request())

			if tt.expectError {
				if handlerCalled {
					t.Error("Handler was called but shouldn't have been")
				}

				if rr.Code != tt.expectedStatus {
					t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
				}

			} else {
				if !handlerCalled {
					t.Error("Handler was not called")
				}
			}
		})
	}
}
