package middleware

import (
	"log"
	"net/http"
	"strings"

	"memtravel/auth"
)

// AuthMiddleware adds a authorization layer for the endpoints that need auth to be accessed
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)

		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		valid, err := auth.VerifyToken(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
