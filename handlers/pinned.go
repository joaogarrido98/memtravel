package handlers

import (
	"fmt"
	"log"
	"memtravel/db"
	"memtravel/middleware"
	"net/http"
	"strconv"
)

func (handler *Handler) AddPinnedHandler(w http.ResponseWriter, r *http.Request) {
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

	tripID, deferredErr := strconv.Atoi(r.PathValue(tripParamID))
	if deferredErr != nil {
		return
	}

	rows, deferredErr := handler.database.Query(db.TripBelongsToUser, userID, tripID)
	if deferredErr != nil {
		return
	}

	defer rows.Close()

	if !rows.Next() {
		deferredErr = fmt.Errorf("%d trip does not belong to user", tripID)
		return
	}

	deferredErr = handler.database.ExecQuery(db.AddPinned, userID, tripID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, "")

}

func (handler *Handler) RemovePinnedHandler(w http.ResponseWriter, r *http.Request) {
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

	tripID, deferredErr := strconv.Atoi(r.PathValue(tripParamID))
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.ExecQuery(db.RemovePinned, userID, tripID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, "")
}
