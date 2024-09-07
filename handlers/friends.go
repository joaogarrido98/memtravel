package handlers

import (
	"errors"
	"fmt"
	"log"
	"memtravel/db"
	"memtravel/middleware"
	"net/http"
	"strconv"
)

const (
	declineFriendRequest = "decline"
	acceptFriendRequest  = "accept"
	removeFriendRequest  = "remove"
	addNewFriendRequest  = "add"
)

var handlerTypes = map[string]struct{}{
	declineFriendRequest: {},
	acceptFriendRequest:  {},
	removeFriendRequest:  {},
	addNewFriendRequest:  {},
}

func (handler *Handler) FriendRequestHandler(w http.ResponseWriter, r *http.Request) {
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

	friendParam := r.URL.Query().Get(friendParamID)
	friendID, deferredErr := strconv.Atoi(friendParam)
	if deferredErr != nil {
		return
	}

	if userID == friendID {
		deferredErr = errors.New("user id and friend id cannot be the same")
		return
	}

	requestType := r.PathValue(friendRequestParamID)
	_, validType := handlerTypes[requestType]
	if !validType {
		deferredErr = fmt.Errorf("%s is not a valid type", requestType)
		return
	}

	switch requestType {
	case addNewFriendRequest:
		row := handler.database.QueryRow(db.CheckIfUserHasSpecificFriend, userID, friendID)

		var exists int
		deferredErr = row.Scan(&exists)
		if deferredErr != nil {
			return
		}

		if exists == 0 {
			deferredErr = handler.database.ExecQuery(db.AddFriendRequest, userID, friendID)
			break
		}

		if row.Err() != nil {
			deferredErr = row.Err()
			return
		}

		deferredErr = fmt.Errorf("%s is already an existing friend", friendParam)
	case acceptFriendRequest:
		deferredErr = handler.database.ExecTransaction(db.Transaction{
			db.RemoveFromFriendsRequest: {friendID, userID},
			db.AddNewFriend:             {friendID, userID},
		})
	case declineFriendRequest:
		deferredErr = handler.database.ExecQuery(db.DeclineFriendRequest, friendID, userID)
	case removeFriendRequest:
		deferredErr = handler.database.ExecQuery(db.RemoveFriendRequest, userID, friendID)
	}

	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, "")
}

func (handler *Handler) RemoveFriendHandler(w http.ResponseWriter, r *http.Request) {
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

	friendParam := r.URL.Query().Get(friendParamID)
	friendID, deferredErr := strconv.Atoi(friendParam)
	if deferredErr != nil {
		return
	}

	deferredErr = handler.database.ExecQuery(db.RemoveFriend, userID, friendID)
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, "")
}

func (handler *Handler) GetFriendsHandler(w http.ResponseWriter, r *http.Request) {
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

	var friends []User

	rows, deferredErr := handler.database.Query(db.GetAllFriends, userID)
	if deferredErr != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var friend User

		deferredErr = rows.Scan(&friend.UserID, &friend.FullName, &friend.ProfilePic)
		if deferredErr != nil {
			return
		}

		friends = append(friends, friend)
	}

	deferredErr = rows.Err()
	if deferredErr != nil {
		return
	}

	deferredErr = writeServerResponse(w, true, friends)
}

func (handler *Handler) SearchFriendsHandler(w http.ResponseWriter, r *http.Request) {
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
