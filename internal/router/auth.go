package router

import (
	"net/http"
	"strings"
)

func (cfg *ApiConfig) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header missing")
			return
		}

		cfg.apiKey = strings.TrimPrefix(authHeader, "ApiKey ")
		dbUser, err := cfg.DB.GetUserByApikey(r.Context(), cfg.apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid ApiKey")
			return
		}

		cfg.user = databaseUserToUser(dbUser)

		next.ServeHTTP(w, r)
	})
}
