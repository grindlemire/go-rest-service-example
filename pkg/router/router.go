package router

// router manages the routes of our server. This can get a lot more complicated but allows us to create arbitrarily
// complex middleware and handlers (for example if either a middleware or handler needed a 3rd party connection to a database
// or complex configuration)

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/grindlemire/go-rest-service-example/pkg/handlers"
	"github.com/grindlemire/go-rest-service-example/pkg/middleware"
)

// NewRouter creates a new mux router with all our handlers configured
func NewRouter() (r *mux.Router, err error) {
	r = mux.NewRouter()
	// fingerprint every request coming into the server. This is always our outermost layer of middleware
	r.Use(middleware.RequestFingerprinter)
	// record metrics for every request. This is always our second outermost layer of middleware
	r.Use(middleware.MetricsRecorder)
	v1Router := r.PathPrefix("/v1").Subrouter()

	// Protected Paths
	// Create a subrouter for our authed routes. Add in the auth middleware
	authRouter := v1Router.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.NewAuthenticator().Authenticate)
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
	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFoundPage).GetHandler()

	return r, nil
}
