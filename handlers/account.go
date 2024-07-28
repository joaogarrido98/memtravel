package handlers

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"unicode"

	"memtravel/auth"
	"memtravel/configs"
	"memtravel/language"
	"memtravel/middleware"
	"memtravel/types"
)

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

	languageID := r.URL.Query().Get("lid")
	if !language.SupportedLanguage(languageID) {
		deferredErr = fmt.Errorf("languageID is not supported: got %s", languageID)
		return
	}

	if strings.TrimSpace(loginRequest.Email) == "" || strings.TrimSpace(loginRequest.Password) == "" {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation(languageID, language.EmptyEmailPassword), nil)
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
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation(languageID, language.PasswordInvalid), nil)
		return
	}

	if !userData.Active {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation(languageID, language.InactiveUser), nil)
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

	languageID := r.URL.Query().Get("lid")
	if !language.SupportedLanguage(languageID) {
		deferredErr = fmt.Errorf("languageID is not supported: got %s", languageID)
		return
	}

	var recoverPasswordRequest types.User

	deferredErr = readBody(r, &recoverPasswordRequest)
	if deferredErr != nil {
		return
	}

	if recoverPasswordRequest.Email == "" {
		deferredErr = fmt.Errorf("email cannot be empty")
		return
	}

	newPassword := createNewPassword()

	hashPassword, deferredErr := auth.HashPassword(newPassword)
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.Update("UPDATE users SET password=$1 WHERE email=$2", hashPassword, recoverPasswordRequest.Email)
	if deferredErr != nil {
		return
	}

	deferredErr = sendEmail(
		[]string{recoverPasswordRequest.Email},
		"./templates/recover.gohtml",
		language.GetTranslation(languageID, language.PasswordRecover),
		types.RecoverPasswordTemplate{
			Password: newPassword,
		},
	)
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

	var passwordChangeRequest types.ChangePassword

	deferredErr = readBody(r, &passwordChangeRequest)
	if deferredErr != nil {
		return
	}

	languageID := r.URL.Query().Get("lid")
	if !language.SupportedLanguage(languageID) {
		deferredErr = fmt.Errorf("languageID is not supported: got %s", languageID)
		return
	}

	if strings.TrimSpace(passwordChangeRequest.NewPassword) == "" || strings.TrimSpace(passwordChangeRequest.OldPassword) == "" {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation(languageID, language.EmptyOldNewPassword), nil)
		return
	}

	if !newPasswordIsValid(passwordChangeRequest.NewPassword) {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation(languageID, language.NewPasswordInvalid), nil)
		return
	}

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
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation(languageID, language.ChagePasswordInvalid), nil)
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

	deferredErr = writeServerResponse(w, "", language.GetTranslation(languageID, language.PasswordChanged), nil)
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

	deferredErr = handler.database.Update("UPDATE users SET active=$1 WHERE id=$2", false, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, "", "", nil)
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

func newPasswordIsValid(password string) bool {
	if len(password) <= 8 || len(password) > 32 {
		return false
	}

	var hasLowerCase = false
	var hasUpperCase = false
	var hasNumber = false
	var hasSpecial = false

	for _, char := range password {
		switch {
		case unicode.IsSpace(char):
			return false
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char), unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasLowerCase && hasUpperCase && hasNumber && hasSpecial
}

func createNewPassword() string {
	var buf bytes.Buffer
	for i := uint(0); i < 32; i++ {
		buf.WriteByte(configs.Envs.PasswordCreation[rand.Intn(len(configs.Envs.PasswordCreation))])
	}

	return buf.String()
}
