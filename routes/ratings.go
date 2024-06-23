package routes

import (
	"memtravel/utils"
	"net/http"
)

func AddRatingHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}
