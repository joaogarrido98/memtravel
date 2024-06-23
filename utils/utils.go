package utils

import (
	"log"
	"net/http"
)

func DeferredHandler(err error, w http.ResponseWriter) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
