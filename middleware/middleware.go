package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"memtravel/auth"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	AuthUserID ContextKey = "context.auth.userID"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// CreateStack creates a middleware that executes all the passed middlewares
func CreateStack(middleware ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}

		return next
	}
}

// LogMiddleware adds a logger middleware to the routes
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Printf("result: [%d] method: [%s], path: [%s] duration: [%s]", wrapped.statusCode, r.Method, r.URL.Path, time.Since(start).String())
	})
}

// AuthMiddleware adds a authorization layer for the endpoints that need auth to be accessed
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		verifiedToken, err := auth.VerifyToken(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !verifiedToken.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, exists := verifiedToken.Claims.(jwt.MapClaims)
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims["user"] == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		request := r.WithContext(context.WithValue(r.Context(), AuthUserID, claims["user"]))

		next.ServeHTTP(w, request)
	})
}
