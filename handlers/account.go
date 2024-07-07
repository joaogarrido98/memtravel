package handlers

import (
	"log"
	"net/http"
	"strings"

	"memtravel/auth"
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

	loginRequest := new(types.User)

	deferredErr = readBody(r, loginRequest)
	if deferredErr != nil {
		return
	}

	if strings.TrimSpace(loginRequest.Email) == "" || strings.TrimSpace(loginRequest.Password) == "" {
		_, deferredErr = w.Write([]byte("invalid"))
		return
	}

	var userData types.User

	rows := handler.database.QueryRow("SELECT email, password, active FROM Users WHERE email = $1", loginRequest.Email)

	deferredErr = rows.Scan(&userData.Email, &userData.Password, &userData.Active)
	if deferredErr != nil {
		return
	}

	passwordValid, deferredErr := auth.CompareHash(loginRequest.Password, userData.Password)
	if deferredErr != nil {
		return
	}

	if !passwordValid || !userData.Active {
		_, deferredErr = w.Write([]byte("invalid"))
		return
	}

	token, deferredErr := auth.CreateToken(userData.UserID)
	if deferredErr != nil {
		return
	}

	_, deferredErr = w.Write([]byte(token))
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
