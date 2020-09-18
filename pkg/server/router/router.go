package router

// router manages the routes of our server. This can get a lot more complicated but allows us to create arbitrarily
// complex middleware and handlers (for example if either a middleware or handler needed a 3rd party connection to a database
// or complex configuration)

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grindlemire/go-rest-service-example/pkg/server/endpoint"
	"github.com/grindlemire/go-rest-service-example/pkg/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter creates a new mux router with all our handlers configured
func NewRouter(authed []endpoint.Endpoint, public []endpoint.Endpoint) (r *mux.Router, err error) {
	r = mux.NewRouter()
	// fingerprint every request coming into the server. This is always our outermost layer of middleware
	r.Use(middleware.RequestFingerprinter)
	// record metrics for every request. This is always our second outermost layer of middleware
	r.Use(middleware.MetricsRecorder)
	v1Router := r.PathPrefix("/v1").Subrouter()

	// Protected Paths
	// Create a subrouter for our authed routes. Add in the auth middleware
	authedRouter := v1Router.PathPrefix("/").Subrouter()
	authedRouter.Use(middleware.NewAuthenticator().Authenticate)
	for _, endpoint := range authed {
		authedRouter.NewRoute().
			Methods(endpoint.Method).
			Path(endpoint.Path).
			HandlerFunc(endpoint.Handler)
	}

	// Public paths
	// Create a subrouter for our public paths
	publicRouter := v1Router.PathPrefix("/").Subrouter()
	for _, endpoint := range public {
		publicRouter.NewRoute().
			Methods(endpoint.Method).
			Path(endpoint.Path).
			HandlerFunc(endpoint.Handler)
	}

	// Metrics endpoint for prometheus
	r.NewRoute().
		Methods(http.MethodGet).
		Path("/metrics").
		Handler(promhttp.Handler())

	// Bare Response redirect to whatever URL we want
	r.NewRoute().
		Path("/").
		Handler(http.RedirectHandler("/v1/public", http.StatusMovedPermanently))

	// This is required because the default NotFoundHandler will bypass all the middleware but we don't want that.
	// Note that this needs to go last in our router so we don't wildcard over the rest of our routes.
	// See https://stackoverflow.com/questions/43613311/make-a-custom-404-with-golang-and-mux
	r.NotFoundHandler = r.NewRoute().HandlerFunc(endpoint.NotFoundPage).GetHandler()

	return r, nil
}
