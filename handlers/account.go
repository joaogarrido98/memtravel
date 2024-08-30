package handlers

import (
	"bytes"
	"math/rand"
	"net/http"
	"strings"
	"unicode"

	"memtravel/auth"
	"memtravel/configs"
	"memtravel/db"
	"memtravel/language"
	"memtravel/middleware"
	"memtravel/types"
)

func (handler *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	var loginRequest types.User

	deferredErr = readBody(r, &loginRequest)
	if deferredErr != nil {
		return
	}

	languageID := r.URL.Query().Get(languageParamID)
	if !language.SupportedLanguage(languageID) {
		deferredErr = errorLanguageID
		return
	}

	if strings.TrimSpace(loginRequest.Email) == "" || strings.TrimSpace(loginRequest.Password) == "" {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation(languageID, language.EmptyEmailPassword), nil)
		return
	}

	var userData types.User

	rows := handler.database.QueryRow(db.GetUserLogin, loginRequest.Email)

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

func (handler *Handler) PasswordRecoverHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	languageID := r.URL.Query().Get(languageParamID)
	if !language.SupportedLanguage(languageID) {
		deferredErr = errorLanguageID
		return
	}

	var recoverPasswordRequest types.User

	deferredErr = readBody(r, &recoverPasswordRequest)
	if deferredErr != nil {
		return
	}

	if recoverPasswordRequest.Email == "" {
		deferredErr = writeServerResponse(w, "invalid", language.GetTranslation(languageID, language.EmptyEmailPassword), nil)
		return
	}

	newPassword := createNewPassword()

	hashPassword, deferredErr := auth.HashPassword(newPassword)
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.Update(db.UpdateUserPassword, hashPassword, recoverPasswordRequest.Email)
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

	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, "", language.GetTranslation(languageID, language.EmptyEmailPassword), nil)
}

func (handler *Handler) PasswordChangeHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
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

	languageID := r.URL.Query().Get(languageParamID)
	if !language.SupportedLanguage(languageID) {
		deferredErr = errorLanguageID
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

	rows := handler.database.QueryRow(db.GetPasswordDetails, userID)

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

	deferredErr = handler.database.Update(db.UpdatePassword, hashedPassword, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, "", language.GetTranslation(languageID, language.PasswordChanged), nil)
}

func (handler *Handler) CloseAccountHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	userID := r.Context().Value(middleware.AuthUserID)

	languageID := r.URL.Query().Get(languageParamID)
	if !language.SupportedLanguage(languageID) {
		deferredErr = errorLanguageID
		return
	}

	deferredErr = handler.database.Update(db.UpdateUserStatus, false, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, "", language.GetTranslation(languageID, language.AccountClose), nil)
}

func (handler *Handler) AccountInformationHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) AccountInformationEditHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func newPasswordIsValid(password string) bool {
	if len(password) < 8 || len(password) > 32 {
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
