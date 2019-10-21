package middleware

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/gorilla/mux"
	"github.com/grindlemire/log"
)

// Authenticator manages the authentication for requests. It is empty here but could be extended
// to pull in state (for example if you wanted to check user auth against
// a persistent store then you would keep the connection in here).
type Authenticator struct {
	// a db connection or something would go here
}

// NewAuthenticator creates a new authenticator struct that could be used to authenticate requests.
func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

// Authenticate authenticates requests for the authenticator
func (a Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the fingerprint of the request for logging and validation
		fingerprint, err := GetRequestFingerprint(r)
		if err != nil {
			log.Errorf("internal error for path [%s]: %v", r.URL.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := a.validateRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			log.Warnf("request[%s] to [%s] from [%s] unauthorized: %v", fingerprint.GetID(), r.URL.String(), fingerprint.GetSource(), err)
			return
		}
		// set the user in the fingerprint for downstream consumption
		fingerprint.SetUser(user)

		next.ServeHTTP(w, r)
	})
}

// validateRequest validates that the request is properly authenticated and authorized for the endpoint
func (a Authenticator) validateRequest(r *http.Request) (user string, err error) {
	// get our route variables out of the path
	pathVars := mux.Vars(r)

	// check the basic auth. In a real system you would use something like jwt auth here
	username, _, found := r.BasicAuth()
	if !found {
		return "", errors.New("unable to find basic auth in request")
	}

	// This random sleep simulates an io operation for checking auth against a db or somthing. It also helps to show the latencies
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	if username != pathVars["id"] {
		return "", errors.Errorf("unable to authenticate [%s]. Invalid credentials", username)
	}

	return username, nil
}
