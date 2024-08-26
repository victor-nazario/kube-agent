package authz

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/victor-nazario/kube-agent/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationHandler(t *testing.T) {
	m := make(map[string]user.User)
	m["operant"] = user.User{
		UserName: "operant",
		Roles:    []string{"cluster-owner"},
	}

	a, err := NewAuthorizer(m)
	if err != nil {
		t.Fatal(err)
	}

	handler := Middleware(a)

	var authenticationTests = []struct {
		code int
		user string
		pass string
		name string
	}{
		{200, "operant", "secret", "authorised"},
		{403, "random", "badpassword", "unauthorised"},
		{403, "", "", "unauthorised empty creds"},
	}

	for _, tt := range authenticationTests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req, err := http.NewRequest("POST", "/query", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.SetBasicAuth(tt.user, tt.pass)

			authFmt := fmt.Sprintf("%s:%s", tt.user, tt.pass)
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authFmt))))

			mockHandlerAuth := func(w http.ResponseWriter, r *http.Request) {
				if username, _, ok := r.BasicAuth(); !ok && username != tt.user {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusOK)
			}

			req = req.WithContext(context.WithValue(req.Context(), "auth", true))

			handler(http.HandlerFunc(mockHandlerAuth)).ServeHTTP(rr, req)

			if rr.Code != tt.code {
				t.Fatalf("Wrong status code returned from Authz handler. Returned code: %d", rr.Code)
			}
		})
	}
}
