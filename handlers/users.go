package handlers

import (
	"log"
	"memtravel/middleware"
	"net/http"
)

type (
	UserPage struct {
		UserID         int          `json:"userid"`
		IsPrivate      bool         `json:"isPrivate"`
		IsFriend       bool         `json:"isFriend"`
		Name           string       `json:"name"`
		ProfilePicture string       `json:"profilepic"`
		Country        string       `json:"country"` // This can be many...
		FriendsSince   string       `json:"friendsSince,omitempty"`
		MemberSince    string       `json:"memberSince"`
		TotalFriends   int          `json:"totalFriends"`
		Bio            string       `json:"bio"`
		PinnedTrips    []PinnedTrip `json:"pinned,omitempty"`
		Stats          []Stats      `json:"stats,omitempty"`
	}

	PinnedTrip struct {
		TripID    int
		Cover     string
		Country   string
		StartDate string
	}

	Stats struct {
		TotalTrips     int
		DaysTravelling int
		TotalCountries int
		TotalCities    int
	}
)

func (handler *Handler) SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
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

}

func (handler *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
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

}

func (handler *Handler) UserEditHandler(w http.ResponseWriter, r *http.Request) {
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

	// TODO
}
