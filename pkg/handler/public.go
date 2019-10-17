package handler

import (
	"net/http"

	"github.com/grindlemire/log"
)

// PublicPage handles the public endpoint response
func PublicPage(w http.ResponseWriter, r *http.Request) {
	log.Infof("the public page handler was hit with route: [%s]", r.URL.String())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("public response works"))
}
