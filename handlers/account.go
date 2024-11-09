package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode"

	"memtravel/auth"
	"memtravel/configs"
	"memtravel/db"
	"memtravel/language"
	"memtravel/middleware"
)

type (
	// ChangePassword is the blueprint for the change password request
	ChangePassword struct {
		OldPassword string `json:"op"`
		NewPassword string `json:"np"`
	}

	// WelcomeTemplate is the blueprint for the new user welcome email
	WelcomeTemplate struct {
		Link string
	}

	// RecoverPasswordTemplate is the blueprint for the reset password email
	RecoverPasswordTemplate struct {
		Password string
	}
)

func (handler *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s]",
				deferredErr.Error(),
				r.Context().Value(middleware.RequestContextID),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	var loginRequest User

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
		deferredErr = errorInvalidRequestData
		return
	}

	var userData User

	row := handler.database.QueryRow(db.GetUserLogin, loginRequest.Email)

	deferredErr = row.Scan(&userData.UserID, &userData.Email, &userData.Password, &userData.Active, &userData.LoginAttempt, &userData.FullName)
	if deferredErr != nil && deferredErr != sql.ErrNoRows {
		return
	}

	if deferredErr != nil && deferredErr == sql.ErrNoRows {
		deferredErr = writeServerResponse(w, false, language.GetTranslation(languageID, language.AccountNotExisting))
		return
	}

	if userData.LoginAttempt >= 5 {
		deferredErr = writeServerResponse(w, false, language.GetTranslation(languageID, language.BlockedLogin))
		return
	}

	if !userData.Active {
		deferredErr = writeServerResponse(w, false, language.GetTranslation(languageID, language.InactiveUser))
		return
	}

	passwordValid, deferredErr := auth.CompareHash(loginRequest.Password, userData.Password)
	if deferredErr != nil {
		return
	}

	if !passwordValid {
		deferredErr = handler.database.ExecQuery(db.UpdateLoginCounter, userData.UserID)
		if deferredErr != nil {
			return
		}

		deferredErr = writeServerResponse(w, false, language.GetTranslation(languageID, language.PasswordInvalid))
		return
	}

	deferredErr = handler.database.ExecQuery(db.ResetLoginCounter, userData.UserID)
	if deferredErr != nil {
		return
	}

	token, deferredErr := auth.CreateToken(strconv.Itoa(userData.UserID))
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, User{Token: token, FullName: userData.FullName})
}

func (handler *Handler) PasswordRecoverHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s]",
				deferredErr.Error(),
				r.Context().Value(middleware.RequestContextID),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	languageID := r.URL.Query().Get(languageParamID)
	if !language.SupportedLanguage(languageID) {
		deferredErr = errorLanguageID
		return
	}

	var recoverPasswordRequest User

	deferredErr = readBody(r, &recoverPasswordRequest)
	if deferredErr != nil {
		return
	}

	if recoverPasswordRequest.Email == "" {
		deferredErr = errorInvalidRequestData
		return
	}

	var emailExists bool
	deferredErr = handler.database.QueryRow(db.EmailExists, recoverPasswordRequest.Email).Scan(&emailExists)
	if deferredErr != nil {
		return
	}

	if !emailExists {
		deferredErr = writeServerResponse(w, true, language.GetTranslation(languageID, language.PasswordRecoverySuccess))
		return
	}

	newPassword := generateRandomString()

	hashPassword, deferredErr := auth.HashPassword(newPassword)
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.ExecQuery(db.UpdateUserPassword, hashPassword, recoverPasswordRequest.Email)
	if deferredErr != nil {
		return
	}

	deferredErr = sendEmail(
		[]string{recoverPasswordRequest.Email},
		"recover.html",
		language.GetTranslation(languageID, language.PasswordRecover),
		RecoverPasswordTemplate{
			Password: newPassword,
		},
		handler.tmpl,
	)

	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, language.GetTranslation(languageID, language.PasswordRecoverySuccess))
}

func (handler *Handler) PasswordChangeHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s], user_id: [%s]",
				deferredErr.Error(),
				r.Context().Value(middleware.RequestContextID),
				r.Context().Value(middleware.AuthUserID),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	userID := r.Context().Value(middleware.AuthUserID)

	var passwordChangeRequest ChangePassword

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
		deferredErr = errorInvalidRequestData
		return
	}

	if !newPasswordIsValid(passwordChangeRequest.NewPassword) {
		deferredErr = errorInvalidRequestData
		return
	}

	var userData User

	row := handler.database.QueryRow(db.GetPasswordDetails, userID)

	deferredErr = row.Scan(&userData.UserID, &userData.Password)
	if deferredErr != nil {
		return
	}

	passwordValid, deferredErr := auth.CompareHash(passwordChangeRequest.OldPassword, userData.Password)
	if deferredErr != nil {
		return
	}

	if !passwordValid {
		deferredErr = writeServerResponse(w, false, language.GetTranslation(languageID, language.ChagePasswordInvalid))
		return
	}

	hashedPassword, deferredErr := auth.HashPassword(passwordChangeRequest.NewPassword)
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.ExecQuery(db.UpdatePassword, hashedPassword, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, "")
}

func (handler *Handler) PrivacyStatusHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s], user_id: [%s]",
				deferredErr.Error(),
				r.Context().Value(middleware.RequestContextID),
				r.Context().Value(middleware.AuthUserID),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	userID := r.Context().Value(middleware.AuthUserID)

	currentPrivacyStatus := r.URL.Query().Get(privacyParamID)
	if currentPrivacyStatus == "" {
		deferredErr = fmt.Errorf("invalid privacy status")
		return
	}

	status, deferredErr := strconv.ParseBool(currentPrivacyStatus)
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.ExecQuery(db.UpdateUserPrivacyStatus, !status, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, "")
}

func (handler *Handler) CloseAccountHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s], user_id: [%s]",
				deferredErr.Error(),
				r.Context().Value(middleware.RequestContextID),
				r.Context().Value(middleware.AuthUserID),
			)
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

	deferredErr = handler.database.ExecQuery(db.UpdateUserActiveStatus, false, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, language.GetTranslation(languageID, language.AccountClose))
}

func (handler *Handler) ActivateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s]",
				deferredErr.Error(),
				r.Context().Value(middleware.RequestContextID),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	code := r.PathValue(codeParamID)
	if len(code) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to activate account"))
		return
	}

	var email string
	var databaseCode string

	row := handler.database.QueryRow(db.GetActivationCode, code)
	deferredErr = row.Scan(&databaseCode, &email)
	if deferredErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to activate account"))
		return
	}

	deferredErr = handler.database.ExecTransaction(
		[]db.Transaction{
			{
				Query:  db.RemoveActivationCode,
				Params: []any{code, email},
			},
			{
				Query:  db.ActivateUser,
				Params: []any{email},
			},
		},
	)

	if deferredErr != nil {
		return
	}

	http.ServeFile(w, r, "./static/activeaccount.html")
}

func (handler *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s]",
				deferredErr.Error(),
				r.Context().Value(middleware.RequestContextID),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	languageID := r.URL.Query().Get(languageParamID)
	if !language.SupportedLanguage(languageID) {
		deferredErr = errorLanguageID
		return
	}

	var registerRequest User

	deferredErr = readBody(r, &registerRequest)
	if deferredErr != nil {
		return
	}

	if strings.TrimSpace(registerRequest.FullName) == "" {
		deferredErr = errorInvalidRequestData
		return
	}

	if len(registerRequest.FullName) >= 45 {
		deferredErr = errorInvalidRequestData
		return
	}

	dateOfBirth, deferredErr := time.Parse(time.DateOnly, registerRequest.DoB)
	if deferredErr != nil {
		return
	}

	cutOffDate := time.Now().AddDate(-16, 0, 0)

	if !cutOffDate.After(dateOfBirth) {
		deferredErr = errorInvalidRequestData
		return
	}

	if !newPasswordIsValid(strings.TrimSpace(registerRequest.Password)) {
		deferredErr = errorInvalidRequestData
		return
	}

	_, deferredErr = mail.ParseAddress(registerRequest.Email)
	if deferredErr != nil {
		return
	}

	rows, deferredErr := handler.database.Query(db.GetUserAccount, registerRequest.Email)
	if deferredErr != nil {
		return
	}

	defer rows.Close()

	if rows.Next() {
		deferredErr = writeServerResponse(w, false, language.GetTranslation(languageID, language.AccountExisting))
		return
	}

	hashedPassword, deferredErr := auth.HashPassword(registerRequest.Password)
	if deferredErr != nil {
		return
	}

	activationCode := generateRandomString()

	deferredErr = handler.database.ExecTransaction(
		[]db.Transaction{
			{
				Query:  db.AddNewUser,
				Params: []any{registerRequest.Email, hashedPassword, registerRequest.FullName, registerRequest.DoB, registerRequest.Country},
			},
			{
				Query:  db.AddUserFlags,
				Params: []any{registerRequest.Email},
			},
			{
				Query:  db.AddUserCounters,
				Params: []any{registerRequest.Email},
			},
			{
				Query:  db.AddActivationCode,
				Params: []any{activationCode, registerRequest.Email},
			},
		},
	)

	if deferredErr != nil {
		return
	}

	deferredErr = sendEmail(
		[]string{registerRequest.Email},
		"welcome.html",
		language.GetTranslation(languageID, language.Welcome),
		WelcomeTemplate{
			Link: "http://localhost:8080/account/activate/" + activationCode, //TODO: CHANGE LINK TO BE REAL ONE
		},
		handler.tmpl,
	)

	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, language.GetTranslation(languageID, language.AccountCreated))
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

func generateRandomString() string {
	var buf bytes.Buffer
	for i := uint(0); i < 32; i++ {
		buf.WriteByte(configs.Envs.RandomCreator[rand.Intn(len(configs.Envs.RandomCreator))])
	}

	return buf.String()
}
