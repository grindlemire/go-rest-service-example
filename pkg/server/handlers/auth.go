package handlers

import (
	"net/http"

	"github.com/grindlemire/log"

	"github.com/grindlemire/go-rest-service-example/pkg/middleware"
)

// AuthedPage handles the user authed endpoint response
func AuthedPage(w http.ResponseWriter, r *http.Request) {
	fingerprint, err := middleware.GetRequestFingerprint(r)
	if err != nil {
		log.Errorf("internal error for path [%s]: %v", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("request[%s] hit authed page [%s]: user: [%s]", fingerprint.GetID(), r.URL.String(), fingerprint.GetUser())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user authed response works"))
}
