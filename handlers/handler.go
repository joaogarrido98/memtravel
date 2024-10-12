package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"net/smtp"

	"memtravel/configs"
	"memtravel/db"
)

type (
	// User is the blueprint for the user data
	User struct {
		UserID         string  `json:"userid"`
		FullName       string  `json:"fullname"`
		Email          string  `json:"email,omitempty"`
		Password       string  `json:"password,omitempty"`
		DoB            string  `json:"dob,omitempty"`
		Bio            *string `json:"bio,omitempty"`
		Country        int     `json:"country,omitempty"` // where the user is originally from
		Active         bool    `json:"active,omitempty"`
		LoginAttempt   int     `json:"loginattempt,omitempty"`
		AccountCreated bool    `json:"accountcreated,omitempty"`
		ProfilePic     *string `json:"profilepic,omitempty"`
	}

	// ServerResponse holds the generic type for all responses in
	ServerResponse struct {
		Status bool        `json:"st"`
		Data   interface{} `json:"dt"`
	}

	// Handler object that holds all needed attributes for the handlers
	Handler struct {
		database db.Database
		tmpl     *template.Template
	}
)

const (
	languageParamID      string = "lid"
	pathParamID          string = "id"
	codeParamID          string = "code"
	friendRequestParamID string = "type"
	friendParamID        string = "friend"
	privacyParamID       string = "pid"
	tripParamID          string = "tpid"
)

var (
	errorLanguageID         = errors.New("languageID is not supported")
	errorPathValueNotFound  = errors.New("path value not found")
	errorInvalidRequestData = errors.New("invalid request data")
)

// NewHandler creates a new object
func NewHandler(db db.Database, tmpl *template.Template) *Handler {
	return &Handler{
		database: db,
		tmpl:     tmpl,
	}
}

func readBody(r *http.Request, into any) error {
	if r.Body == nil {
		return errors.New("request body cannot be empty")
	}

	err := json.NewDecoder(r.Body).Decode(into)
	if err != nil {
		return err
	}

	return nil
}

func writeServerResponse(w http.ResponseWriter, status bool, data interface{}) error {
	serverResponse := ServerResponse{
		Status: status,
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(serverResponse)
}

func sendEmail(sendTo []string, emailType string, subject string, context any, t *template.Template) error {
	auth := smtp.PlainAuth("", configs.Envs.EmailFrom, configs.Envs.EmailPassword, configs.Envs.SMTPHost)

	var body bytes.Buffer
	mimeHeaders := "\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	_, err := body.Write([]byte("Subject: " + subject + mimeHeaders))
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(&body, emailType, context)
	if err != nil {
		return err
	}

	return smtp.SendMail(configs.Envs.SMTPHost+":"+configs.Envs.SMTPPort, auth, configs.Envs.EmailFrom, sendTo, body.Bytes())
}
