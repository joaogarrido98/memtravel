package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"memtravel/configs"
	"net/http"
	"net/smtp"
	"text/template"
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

func sendEmail(sendTo []string, emailType string, subject string, context any) error {
	auth := smtp.PlainAuth("", configs.Envs.EmailFrom, configs.Envs.EmailPassword, configs.Envs.SMTPHost)

	t, err := template.ParseFiles(emailType)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	_, err = body.Write([]byte(fmt.Sprintf("%s\n%s\n\n", subject, mimeHeaders)))
	if err != nil {
		return err
	}

	err = t.Execute(&body, context)
	if err != nil {
		return err
	}

	return smtp.SendMail(configs.Envs.SMTPHost+":"+configs.Envs.SMTPPort, auth, configs.Envs.EmailFrom, sendTo, body.Bytes())
}
