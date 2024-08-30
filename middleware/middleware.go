package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"memtravel/auth"
	"memtravel/types"

	"github.com/golang-jwt/jwt/v5"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

const AuthUserID types.ContextKey = "context.auth.userID"

// CreateStack creates a middleware that executes all the passed middlewares
func CreateStack(middleware ...types.Middleware) types.Middleware {
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

		wrapped := &types.WrappedWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		if wrapped.StatusCode == http.StatusOK {
			logger.Info("", "result", "200", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
			return
		}

		logger.Error("", "result", strconv.Itoa(wrapped.StatusCode), "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
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
