package middleware

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"memtravel/auth"
	"memtravel/ratelimiter"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

const (
	AuthUserID       ContextKey = "context.auth.userID"
	RequestContextID ContextKey = "context.request.id"
)

type (
	// Middleware is the blueprint for the handlerfunc
	Middleware func(next http.HandlerFunc) http.HandlerFunc

	// ContextKey is the blueprint for the request context
	ContextKey string

	// WrappedWriter extends the http.ResponseWriter
	WrappedWriter struct {
		http.ResponseWriter
		StatusCode int
	}
)

// WriteHeader is an extension of the http.ResponseWriter WriteHeader
func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
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

// BaseMiddleware adds a base middleware to the routes
// this middleware makes sure to deal with cors policy, basic logging and rate limiting
func BaseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		start := time.Now()

		contextID := uuid.NewString()
		r = r.WithContext(context.WithValue(r.Context(), RequestContextID, contextID))

		clientIP := getClientIP(r)

		if !ratelimiter.GetGlobalLimiter().Allow(clientIP) {
			contextID, _ := r.Context().Value(RequestContextID).(string)

			logger.Warn(
				contextID,
				"message", "Rate limit exceeded",
				"client_ip", clientIP,
				"method", r.Method,
				"path", r.URL.Path,
			)

			w.Header().Set("Retry-After", "60")
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		wrapped := &WrappedWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		if wrapped.StatusCode != http.StatusOK {
			logger.Error(
				contextID,
				"result", strconv.Itoa(wrapped.StatusCode),
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start).String(),
				"host", r.Host,
				"remoteAddress", r.RemoteAddr,
			)
		}
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

// getClientIP extracts the client IP address from request
// Handles cases where the request may be behind a proxy
func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
