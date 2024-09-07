package handlers

import (
	"log"
	"net/http"
	"strconv"

	"memtravel/db"
	"memtravel/language"
	"memtravel/middleware"
)

func (handler *Handler) AddTripHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) GetUpcomingTripsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) GetPreviousTripsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) EditTripHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func (handler *Handler) RemoveTripHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s], user_id: [%s]", deferredErr.Error(), r.Context().Value(middleware.RequestContextID), r.Context().Value(middleware.AuthUserID))
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

	tripID := r.PathValue("id")
	if tripID == "" {
		deferredErr = errorPathValueNotFound
		return
	}

	tripIDInt, deferredErr := strconv.Atoi(tripID)
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.ExecQuery(db.RemoveTrip, tripIDInt, userID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, nil)
}
