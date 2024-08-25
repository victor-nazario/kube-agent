package authz

import (
	"github.com/victor-nazario/kube-agent/internal/user"
	"log"
	"net/http"
)

const defaultResource = "cluster"

type Authorizer interface {
	HasPermission(userName, action, resource string) bool
}

type authorizer struct {
	users map[string]user.User
	roles Roles
}

func NewAuthorizer(users map[string]user.User) (Authorizer, error) {
	r, err := LoadRoles()
	if err != nil {
		return nil, err
	}

	return authorizer{
		users: users,
		roles: r,
	}, nil
}

// HasPermission validates if a given user has permission to perform an action on a resource
func (a authorizer) HasPermission(userName, action, resource string) bool {
	if usr, ok := a.users[userName]; ok {
		for _, roleName := range usr.Roles {
			role := a.roles[roleName]
			if role == nil {
				return false
			}

			if allow, ok := role[action]; ok {
				for _, act := range allow {
					if act == resource {
						return true
					}
				}
			}
		}
	}
	return false
}

func Middleware(a Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, _, ok := r.BasicAuth()
			action := ActionFromMethod(r)
			if !ok || !a.HasPermission(username, action, defaultResource) {
				log.Printf("User '%s' is denied '%s' on resource '%s'", username, action, defaultResource)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ActionFromMethod(r *http.Request) string {
	switch r.Method {
	case http.MethodGet:
		return "can_read"
	case http.MethodPost, http.MethodPut:
		return "can_write"
	case http.MethodDelete:
		return "can_delete"
	default:
		return ""
	}
}
