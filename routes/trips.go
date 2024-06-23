package routes

import (
	"memtravel/utils"
	"net/http"
)

// Gets a specific trip
func GetTripsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func AddTripHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func GetUpcomingTripsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func GetPreviousTripsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func EditTripHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func RemoveTripHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}
