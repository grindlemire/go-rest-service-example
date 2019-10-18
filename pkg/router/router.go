package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/grindlemire/go-rest-service-example/pkg/config"
	"github.com/grindlemire/go-rest-service-example/pkg/handlers"
	"github.com/grindlemire/go-rest-service-example/pkg/middleware"
)

// NewRouter creates a new mux router with all our handlers configured
func NewRouter(opts config.Opts) (r *mux.Router, err error) {
	r = mux.NewRouter()
	// fingerprint every request coming into the server. This is always our outermost layer of middleware
	r.Use(middleware.RequestFingerprinter)
	// record metrics for every request. This is always our second outermost layer of middleware
	r.Use(middleware.MetricsRecorder)

	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundPage)
	v1Router := r.PathPrefix("/v1").Subrouter()

	// Protected Paths
	// Create a subrouter for our authed routes. Add in the auth middleware
	authRouter := v1Router.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.Authenticator)
	authRouter.NewRoute().
		Methods(http.MethodGet).
		Path("/user/{id:[a-zA-Z0-9]+}").
		HandlerFunc(handlers.AuthedPage)

	// Public paths
	// Create a subrouter for our public paths
	publicRouter := v1Router.PathPrefix("/").Subrouter()
	publicRouter.NewRoute().
		Methods(http.MethodGet).
		Path("/public").
		HandlerFunc(handlers.PublicPage)

	return r, nil
}
