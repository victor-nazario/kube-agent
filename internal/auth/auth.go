package auth

import (
	"context"
	"crypto/subtle"
	"net/http"
)

const (
	// In a real world app we would use hashing, salt and pepper to store user and password information.
	user     = "operant"
	password = "secret"
)

// Authentication provides a handler that authenticates client provided login information.
// Returns a handler that contains in the request information the authentication information for the user.
func Authentication() func(http.Handler) http.Handler {
	// This is an oversimplification of authentication in an API. A production implementation
	// would need a more robust mechanism.

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			var ctx context.Context

			usr, pass, ok := request.BasicAuth()
			if !ok || subtle.ConstantTimeCompare([]byte(usr), []byte(user)) != 1 ||
				subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				ctx = context.WithValue(request.Context(), "auth", false)
			} else {
				ctx = context.WithValue(request.Context(), "auth", true)
			}

			// we set the context with the value here, this later gets used in RBAC to first make sure the user
			// is authenticated
			request = request.WithContext(ctx)
			handler.ServeHTTP(writer, request)
		})
	}
}

// IsUserAuthenticated returns the value associated with the auth key on the given context
func IsUserAuthenticated(ctx context.Context) bool {
	auth := ctx.Value("auth").(bool)
	return auth
}
