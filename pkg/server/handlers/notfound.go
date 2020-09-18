package handlers

import (
	"net/http"

	"github.com/grindlemire/log"

	"github.com/grindlemire/go-rest-service-example/pkg/middleware"
)

// NotFoundPage handles requests where the route was not found
func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	fingerprint, err := middleware.GetRequestFingerprint(r)
	if err != nil {
		log.Errorf("internal error for path [%s]: %v", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("request[%s] hit not found page [%s]", fingerprint.GetID(), r.URL.String())
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("page not found"))
}
