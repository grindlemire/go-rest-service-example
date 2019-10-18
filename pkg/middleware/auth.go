package middleware

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grindlemire/log"
)

// Authenticator authenticates requests
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get our route variables out of the path
		pathVars := mux.Vars(r)

		// get the fingerprint of the request for logging and validation
		fingerprint, err := GetRequestFingerprint(r)
		if err != nil {
			log.Errorf("Unable to get request fingerprint: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// check the basic auth. In a real system you would ues something like jwt auth here
		username, _, found := r.BasicAuth()

		// This random sleep simulates an io operation for checking auth against a db or somthing. It also helps to show the latencies
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		if !found || username != pathVars["id"] {
			w.WriteHeader(http.StatusUnauthorized)
			log.Warnf("request[%s] to [%s] from [%s] unauthorized", fingerprint.GetID(), r.URL.String(), fingerprint.GetSource())
			return
		}

		// set the user in the fingerprint for downstream consumption
		fingerprint.SetUser(username)

		next.ServeHTTP(w, r)
	})
}
