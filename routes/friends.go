package routes

import (
	"memtravel/utils"
	"net/http"
)

func FriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)
}

func AcceptFriendHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func RemoveFriendHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func GetFriendsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func GetFriendHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func SearchFriendsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}
