package handlers

import (
	"net/http"

	"github.com/grindlemire/log"

	"github.com/grindlemire/go-rest-service-example/pkg/middleware"
)

// PublicPage handles the public endpoint response
func PublicPage(w http.ResponseWriter, r *http.Request) {
	fingerprint, err := middleware.GetRequestFingerprint(r)
	if err != nil {
		log.Errorf("internal error for path [%s]: %v", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("request[%s] hit public page [%s]", fingerprint.GetID(), r.URL.String())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("public response works"))
}
