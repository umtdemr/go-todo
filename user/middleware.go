package user

import (
	"github.com/umtdemr/go-todo/server"
	"net/http"
	"strings"
)

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			server.RespondWithError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			server.RespondWithError(w, "please provide Bearer token", http.StatusUnauthorized)
			return
		}

		jwtString := strings.TrimPrefix(tokenString, "Bearer ")

		isTokenValid, validationErr := ValidateJWT(jwtString)

		if validationErr != nil {
			server.RespondWithError(w, validationErr.Error(), http.StatusUnauthorized)
			return
		}

		if !isTokenValid {
			server.RespondWithError(w, "token is not valid", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
