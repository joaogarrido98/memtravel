package types

import (
	"net/http"
)

// ServerResponse holds the generic type for all responses in
type ServerResponse struct {
	Status bool        `json:"st"`
	Data   interface{} `json:"dt,omitempty"`
}

// User is the blueprint for the user data
type User struct {
	UserID   string `json:"id,omitempty"`
	FullName string `json:"fullname,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	DoB      string `json:"dob,omitempty"`
	Country  string `json:"country,omitempty"` // where the user is originally from
	Active   bool   `json:"active,omitempty"`
}

// ChangePassword is the blueprint for the change password request
type ChangePassword struct {
	OldPassword string `json:"op,omitempty"`
	NewPassword string `json:"np,omitempty"`
}

// RecoverPasswordTemplate is the blueprint for the reset password email
type RecoverPasswordTemplate struct {
	Password string
}

// WelcomeTemplate is the blueprint for the new user welcome email
type WelcomeTemplate struct {
	Link string
}

// Middleware is the blueprint for the handlerfunc
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// WrappedWriter extends the http.ResponseWriter
type WrappedWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader is an extension of the http.ResponseWriter WriteHeader
func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

// ContextKey is the blueprint for the request context
type ContextKey string

type Transaction map[string][]any
