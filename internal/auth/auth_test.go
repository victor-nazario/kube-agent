package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsUserAuthenticated(t *testing.T) {
	ctx := context.WithValue(context.Background(), "auth", true)

	if !IsUserAuthenticated(ctx) {
		t.Fatalf("Call to IsUserAuthenticated with authorized context returned false instead of true.")
	}
}

func TestIsUserAuthenticatedShouldReturnFalse(t *testing.T) {
	ctx := context.WithValue(context.Background(), "auth", false)

	if IsUserAuthenticated(ctx) {
		t.Fatalf("Call to IsUserAuthenticated with unauthorized context returned true instead of false.")
	}
}

func TestAuthenticationHandler(t *testing.T) {
	handler := Authentication()

	var authenticationTests = []struct {
		code int
		user string
		pass string
		name string
	}{
		{200, "operant", "secret", "successful login"},
		{401, "random", "badpassword", "invalid username or password"},
		{401, "", "", "empty credentials"},
	}

	for _, tt := range authenticationTests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req, err := http.NewRequest("POST", "/query", nil)
			if err != nil {
				t.Fatal(err)
			}

			authFmt := fmt.Sprintf("%s:%s", tt.user, tt.pass)
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authFmt))))

			mockHandler := func(w http.ResponseWriter, r *http.Request) {
				if !IsUserAuthenticated(r.Context()) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusOK)
			}

			handler(http.HandlerFunc(mockHandler)).ServeHTTP(rr, req)

			if rr.Code != tt.code {
				t.Fatalf("Wrong status code returned from Authentication handler. Returned code: %d", rr.Code)
			}
		})
	}

}
