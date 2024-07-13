package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"memtravel/configs"
	"memtravel/db"
	"memtravel/types"
	"net/http"
	"net/smtp"
	"text/template"
)

// Handler object that holds all needed attributes for the handlers
type Handler struct {
	database db.Database
}

// NewHandler creates a new object
func NewHandler(db db.Database) *Handler {
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

func writeServerResponse(w http.ResponseWriter, response, message string, data interface{}) error {
	serverResponse := types.ServerResponse{
		Response: response,
		Message:  message,
		Data:     data,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(serverResponse)
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
