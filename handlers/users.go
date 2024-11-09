package handlers

import (
	"errors"
	"fmt"
	"log"
	"memtravel/cache"
	"memtravel/middleware"
	"net/http"
	"time"
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
	if len(searchQuery) == 0 {
		deferredErr = errors.New("search query parameter is required")
		return
	}

	if len(searchQuery) < 2 {
		writeServerResponse(w, true, []User{})
	}

	if cachedSearch, ok := userCache.Get(searchQuery); ok {
		writeServerResponse(w, true, cachedSearch)
		return
	}

	query := `
        SELECT userid, fullname, profilepic
        FROM users
        WHERE active = true AND userid != $1 AND fullname LIKE %$2%`

	rows, err := handler.database.Query(query, userID, searchQuery)
	if err != nil {
		deferredErr = fmt.Errorf("failed to query users: %v", err)
		return
	}

	defer rows.Close()

	var results []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.FullName, &user.ProfilePicture); err != nil {
			deferredErr = fmt.Errorf("failed to scan user row: %v", err)
			return
		}
		results = append(results, user)
	}

	if err := rows.Err(); err != nil {
		deferredErr = fmt.Errorf("error iterating over result rows: %v", err)
		return
	}

	userCache.Set(searchQuery, results, 2*time.Minute)

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
