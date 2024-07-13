package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"memtravel/auth"
	"memtravel/language"
	"memtravel/middleware"
	"memtravel/types"
)

// todo: add language id

func (handler *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	var loginRequest types.User

	deferredErr = readBody(r, &loginRequest)
	if deferredErr != nil {
		return
	}

	if strings.TrimSpace(loginRequest.Email) == "" || strings.TrimSpace(loginRequest.Password) == "" {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation("1", language.EmptyEmailPassword), nil)
		return
	}

	var userData types.User

	rows := handler.database.QueryRow("SELECT id, email, password, active FROM Users WHERE email = $1", loginRequest.Email)

	deferredErr = rows.Scan(&userData.UserID, &userData.Email, &userData.Password, &userData.Active)
	if deferredErr != nil {
		return
	}

	passwordValid, deferredErr := auth.CompareHash(loginRequest.Password, userData.Password)
	if deferredErr != nil {
		return
	}

	if !passwordValid {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation("1", language.PasswordInvalid), nil)
		return
	}

	if !userData.Active {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation("1", language.EmptyEmailPassword), nil)
		return
	}

	token, deferredErr := auth.CreateToken(userData.UserID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, "", token, nil)
}

func (handler *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) PasswordRecoverHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) PasswordChangeHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	userID := r.Context().Value(middleware.AuthUserID)
	if userID == "" {
		deferredErr = fmt.Errorf("userID cannot be empty")
		return
	}

	var passwordChangeRequest types.ChangePassword

	deferredErr = readBody(r, &passwordChangeRequest)
	if deferredErr != nil {
		return
	}

	if strings.TrimSpace(passwordChangeRequest.NewPassword) == "" || strings.TrimSpace(passwordChangeRequest.OldPassword) == "" {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation("1", language.EmptyEmailPassword), nil)
		return
	}

	// if !newPasswordIsValid(passwordChangeRequest.NewPassword) {
	// 	deferredErr = writeServerResponse(w, "invalid", language.GetTranslation("1", language.EmptyEmailPassword), nil)
	// 	return
	// }

	var userData types.User

	rows := handler.database.QueryRow("SELECT id, password FROM Users WHERE id = $1", userID)

	deferredErr = rows.Scan(&userData.UserID, &userData.Password)
	if deferredErr != nil {
		return
	}

	passwordValid, deferredErr := auth.CompareHash(passwordChangeRequest.OldPassword, userData.Password)
	if deferredErr != nil {
		return
	}

	if !passwordValid {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation("1", language.PasswordInvalid), nil)
		return
	}

	hashedPassword, deferredErr := auth.HashPassword(passwordChangeRequest.NewPassword)
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.Update("UPDATE users SET password=$1 WHERE id=$2", hashedPassword, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, "", language.GetTranslation("1", language.PasswordChanged), nil)
}

func (handler *Handler) CloseAccountHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	userID := r.Context().Value(middleware.AuthUserID)
	if userID == "" {
		deferredErr = fmt.Errorf("userID cannot be empty")
		return
	}

	deferredErr = handler.database.Update("UPDATE users SET active=$1 WHERE id=$2", false, userID)
}

func (handler *Handler) AccountInformationHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) AccountInformationEditHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}
