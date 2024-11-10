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
		UserID         int          `json:"userid,omitempty"`
		Username       string       `json:"username,omitempty"`
		Email          string       `json:"email,omitempty"`
		Token          string       `json:"token,omitempty"`
		Password       string       `json:"password,omitempty"`
		Active         bool         `json:"active,omitempty"`
		DoB            string       `json:"dob,omitempty"`
		IsPrivate      bool         `json:"isPrivate,omitempty"`
		IsFriend       bool         `json:"isFriend,omitempty"`
		FullName       string       `json:"fullname,omitempty"`
		ProfilePicture string       `json:"profilepic,omitempty"`
		Country        int          `json:"country,omitempty"`
		FriendsSince   string       `json:"friendsSince,omitempty"`
		TotalFriends   int          `json:"totalFriends,omitempty"`
		Bio            string       `json:"bio,omitempty"`
		PinnedTrips    []PinnedTrip `json:"pinned,omitempty"`
		Stats          []Stats      `json:"stats,omitempty"`
		LoginAttempt   int          `json:"loginattempt,omitempty"`
		AccountCreated bool         `json:"accountcreated,omitempty"`
	}

	// ServerResponse holds the generic type for all responses in the api
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
	countryParamID       string = "cid"
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
