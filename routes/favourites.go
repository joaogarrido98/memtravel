package routes

import (
	"memtravel/utils"
	"net/http"
)

func AddFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func RemoveFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}
