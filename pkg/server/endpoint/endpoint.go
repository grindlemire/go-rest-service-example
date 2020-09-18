package endpoint

import (
	"fmt"
	"net/http"

	"github.com/grindlemire/go-rest-service-example/pkg/server/middleware"
	"github.com/grindlemire/log"
)

// Endpoint ...
type Endpoint struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

// CreateEndpoints creates and assembles all the endpoints
func CreateEndpoints() (authedEndpoints []Endpoint, publicEndpoints []Endpoint) {
	authedEndpoints = []Endpoint{
		authed(),
	}

	publicEndpoints = []Endpoint{
		home(),
	}
	return authedEndpoints, publicEndpoints
}

func authed() Endpoint {
	return Endpoint{
		Method: http.MethodGet,
		Path:   "/authed",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fingerprint, err := middleware.GetRequestFingerprint(r)
			if err != nil {
				log.Errorf("internal error for path [%s]: %v", r.URL.Path, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Infof("request[%s] hit authed page [%s]", fingerprint.GetID(), r.URL.String())
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("authed response works"))
		},
	}
}

func home() Endpoint {
	return Endpoint{
		Method: http.MethodGet,
		Path:   "/public",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fingerprint, err := middleware.GetRequestFingerprint(r)
			if err != nil {
				log.Errorf("internal error for path [%s]: %v", r.URL.Path, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Infof("request[%s] hit public page [%s]", fingerprint.GetID(), r.URL.String())
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("public response works"))
		},
	}
}

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

func respond(w http.ResponseWriter, status int, formattedStr string, args ...interface{}) {
	if status != http.StatusOK && status != http.StatusNotFound {
		log.Errorf(formattedStr, args...)
	}
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(formattedStr, args...)))
}
