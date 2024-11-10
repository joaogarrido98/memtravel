package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"memtravel/cache"
	"memtravel/db"
	"memtravel/middleware"
)

type (
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

	SearchUser struct {
		Username       string `json:"username"`
		FullName       string `json:"fullname"`
		ProfilePicture string `json:"profilepicture"`
		Page           int    `json:"page"`
	}
)

var userCache = cache.NewCache()

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

	userID := r.Context().Value(middleware.AuthUserID)

	searchQuery := r.URL.Query().Get("query")
	if len(searchQuery) < 2 {
		writeServerResponse(w, true, []User{})
		return
	}

	page := 1
	limit := 20

	p := r.URL.Query().Get("page")
	if p != "" {
		page, deferredErr = strconv.Atoi(p)
		if deferredErr != nil {
			return
		}
	}

	offset := (page - 1) * limit

	if page == 1 {
		if cachedSearch, ok := userCache.Get(searchQuery); ok {
			writeServerResponse(w, true, cachedSearch)
			return
		}
	}

	rows, deferredErr := handler.database.Query(db.SearchUser, userID, searchQuery, offset)
	if deferredErr != nil {
		deferredErr = fmt.Errorf("failed to query users: %v", deferredErr)
		return
	}

	defer rows.Close()

	var results []SearchUser

	for rows.Next() {
		var user SearchUser

		deferredErr = rows.Scan(&user.FullName, &user.Username, &user.ProfilePicture)
		if deferredErr != nil {
			deferredErr = fmt.Errorf("failed to scan user row: %v", deferredErr)
			return
		}

		user.Page = page

		results = append(results, user)
	}

	deferredErr = rows.Err()
	if deferredErr != nil {
		deferredErr = fmt.Errorf("error iterating over result rows: %v", deferredErr)
		return
	}

	if page == 1 {
		userCache.Set(searchQuery, results, 2*time.Minute)
	}

	deferredErr = writeServerResponse(w, true, results)
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
