package routes

import (
	"memtravel/utils"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func CloseAccountHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func AccountInformationHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}

func AccountInformationEditHandler(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer utils.DeferredHandler(deferredErr, w)

}
