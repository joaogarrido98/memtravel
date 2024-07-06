package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler object that holds all needed attributes for the handlers
type Handler struct {
	database *sql.DB
}

// NewHandler creates a new object
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		database: db,
	}
}

func readBody(r *http.Request, into any) error {
	if r.Body == nil {
		return fmt.Errorf("request body cannot be empty")
	}

	err := json.NewDecoder(r.Body).Decode(into)
	if err != nil {
		return err
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}
