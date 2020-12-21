package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
)

// Predefined roles
const (
	manager = "MANAGER"
	admin   = "ADMIN"
)

// Predefined errors
var (
	ErrNoAuthentication = errors.New("No authentication")
)
var authenticationContextKey = &contextKey{"authentication context"}

// contextKey is struct
type contextKey struct {
	name string
}
// HasAnyRoleFunc is type
type HasAnyRoleFunc func(ctx context.Context, roles ...string) bool

func (c *contextKey) String() string {
	return c.name
}

// IDFunc is 
type IDFunc func(ctx context.Context, token string) (int64, error)

// Authenticate is function
func Authenticate(idFunc IDFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authorization")

			id, err := idFunc(request.Context(), token)
			if err != nil {
				log.Print(err, "Authhhththth")
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(request.Context(), authenticationContextKey, id)
			request = request.WithContext(ctx)

			handler.ServeHTTP(writer, request)
		})
	}
}

// Authentication is function
func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}