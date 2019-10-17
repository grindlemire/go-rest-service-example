package handler

import (
	"net/http"

	"github.com/grindlemire/log"
)

// NotFoundPage handles requests where the route was not found
func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	log.Infof("route not found: [%s]", r.URL.String())
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("page not found"))
}
