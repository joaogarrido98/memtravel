package handlers

import (
	"log"
	"net/http"
)

func FriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
}

func AcceptFriendHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func RemoveFriendHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func GetFriendsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func GetFriendHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}

func SearchFriendsHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: %s", deferredErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

}
