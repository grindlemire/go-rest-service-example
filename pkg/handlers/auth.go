package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/grindlemire/log"
)

// AuthedPage handles the user authed endpoint response
func AuthedPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Infof("the authed page handler was hit with route: [%s]: user: [%s]", r.URL.String(), vars["id"])
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user authed response works"))
}

// NotAuthedPage handles the unauthed endpoint response
func NotAuthedPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("not authorized"))
}
