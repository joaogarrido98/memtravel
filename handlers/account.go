package handlers

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"strings"
	"time"
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
			log.Printf("Error: [%s], context_id: [%s]", deferredErr.Error(), r.Context().Value(middleware.RequestContextID))
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
		deferredErr = errorInvalidRequestData
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
		deferredErr = writeServerResponse(w, false, language.GetTranslation(languageID, language.PasswordInvalid))
		return
	}

	if !userData.Active {
		deferredErr = writeServerResponse(w, false, language.GetTranslation(languageID, language.InactiveUser))
		return
	}

	token, deferredErr := auth.CreateToken(userData.UserID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, token)
}

func (handler *Handler) PasswordRecoverHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s]", deferredErr.Error(), r.Context().Value(middleware.RequestContextID))
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
		deferredErr = errorInvalidRequestData
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
		"./templates/recover.gohtml",
		language.GetTranslation(languageID, language.PasswordRecover),
		types.RecoverPasswordTemplate{
			Password: newPassword,
		},
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
			log.Printf("Error: [%s], context_id: [%s]", deferredErr.Error(), r.Context().Value(middleware.RequestContextID))
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
		deferredErr = errorInvalidRequestData
		return
	}

	if !newPasswordIsValid(passwordChangeRequest.NewPassword) {
		deferredErr = errorInvalidRequestData
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

	deferredErr = writeServerResponse(w, true, language.GetTranslation(languageID, language.PasswordChanged))
}

func (handler *Handler) CloseAccountHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s]", deferredErr.Error(), r.Context().Value(middleware.RequestContextID))
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

	deferredErr = handler.database.ExecQuery(db.UpdateUserStatus, false, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, language.GetTranslation(languageID, language.AccountClose))
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

func (handler *Handler) ActivateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s]", deferredErr.Error(), r.Context().Value(middleware.RequestContextID))
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

	rows := handler.database.QueryRow(db.GetActivationCode, code)
	deferredErr = rows.Scan(&databaseCode, &email)
	if deferredErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to activate account"))
		return
	}

	deferredErr = handler.database.ExecTransaction(
		types.Transaction{
			db.RemoveActivationCode: {code, email},
			db.ActivateUser:         {email},
		},
	)

	if deferredErr != nil {
		return
	}

	w.Write([]byte("\nAccount is now active. Happy Trips."))
}

func (handler *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s]", deferredErr.Error(), r.Context().Value(middleware.RequestContextID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	languageID := r.URL.Query().Get(languageParamID)
	if !language.SupportedLanguage(languageID) {
		deferredErr = errorLanguageID
		return
	}

	var registerRequest types.User

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

	if strings.TrimSpace(registerRequest.Country) == "" {
		deferredErr = errorInvalidRequestData
		return
	}

	if len(registerRequest.Country) >= 25 {
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

	rows, deferredErr := handler.database.Query(db.GetUserLogin, registerRequest.Email)
	if deferredErr != nil {
		return
	}

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
		types.Transaction{
			db.AddNewUser:        {registerRequest.Email, hashedPassword, registerRequest.FullName, registerRequest.DoB, registerRequest.Country},
			db.AddActivationCode: {activationCode, registerRequest.Email},
		},
	)

	if deferredErr != nil {
		return
	}

	deferredErr = sendEmail(
		[]string{registerRequest.Email},
		"./templates/welcome.gohtml",
		language.GetTranslation(languageID, language.Welcome),
		types.WelcomeTemplate{
			Link: "http://localhost:8080/account/activate/" + activationCode, //TODO: CHANGE LINK TO BE REAL ONE
		},
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
