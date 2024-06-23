package handlers

import (
	"encoding/json"
	"log"
	"memtravel/auth"
	"memtravel/types"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	if r.Body == nil {
		deferredErr = errEmptyBody
		return
	}

	var loginRequest types.User

	deferredErr = json.NewDecoder(r.Body).Decode(&loginRequest)
	if deferredErr != nil {
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		_, deferredErr = w.Write([]byte("invalid"))
		return
	}

	var userData types.User
	// TODO: get user password for specific email given from db

	passwordValid, deferredErr := compareHash(loginRequest.Password, userData.Password)
	if deferredErr != nil {
		return
	}

	if !passwordValid {
		_, deferredErr = w.Write([]byte("invalid"))
		return
	}

	token, deferredErr := auth.CreateToken(userData.UserID)
	if deferredErr != nil {
		return
	}

	_, deferredErr = w.Write([]byte(token))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func PasswordRecoverHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func PasswordChangeHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func CloseAccountHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func AccountInformationHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func AccountInformationEditHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}
