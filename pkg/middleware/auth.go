package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/grindlemire/go-rest-service-example/pkg/handler"

	"github.com/grindlemire/log"
)

// Authenticator authenticates requests
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("the authenticator middleware was hit with route: [%s]", r.URL.String())
		vars := mux.Vars(r)
		username, _, found := r.BasicAuth()
		// This is where the authentication would go to make sure users are authorized
		if !found || username != vars["id"] {
			http.HandlerFunc(handler.NotAuthedPage).ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
