package user

import (
	"context"
	"github.com/umtdemr/go-todo/server"
	"net/http"
	"strings"
)

// AuthMiddleware is a middleware that checks if the user is authenticated
// It checks the Authorization header for a Bearer token
// If the token is valid, it adds the user to the context
// If the token is not valid, it returns an error with status code 401
func (service *Service) AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from the Authorization header
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			server.RespondWithError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// check if the token is a Bearer token
		if !strings.HasPrefix(tokenString, "Bearer ") {
			server.RespondWithError(w, "please provide Bearer token", http.StatusUnauthorized)
			return
		}

		// get the JWT string
		jwtString := strings.TrimPrefix(tokenString, "Bearer ")

		// validate the JWT and check if it is valid
		username, validationErr := ValidateJWT(jwtString)

		if validationErr != nil {
			server.RespondWithError(w, validationErr.Error(), http.StatusUnauthorized)
			return
		}

		if username == "" {
			server.RespondWithError(w, "token is not valid", http.StatusUnauthorized)
			return
		}

		// get the user from the database and add it to the context
		user := service.repository.GetUserByUsername(username)
		ctx := context.WithValue(r.Context(), "user", user)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
