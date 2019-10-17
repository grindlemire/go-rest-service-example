package router

import (
	"net/http"

	"github.com/grindlemire/go-rest-service-example/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/grindlemire/go-rest-service-example/pkg/config"
	"github.com/grindlemire/go-rest-service-example/pkg/handler"
)

// NewRouter creates a new mux router with all our handlers configured
func NewRouter(opts config.Opts) (r *mux.Router, err error) {
	r = mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handler.NotFoundPage)
	v1Router := r.PathPrefix("/v1").Subrouter()

	// Protected Paths
	// Create a subrouter for our authed routes. Add in the auth middleware
	authRouter := v1Router.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.Authenticator)
	authRouter.NewRoute().
		Methods(http.MethodGet).
		Path("/user/{id:[a-zA-Z0-9]+}").
		HandlerFunc(handler.AuthedPage)

	// Public paths
	// Create a subrouter for our public paths
	publicRouter := v1Router.PathPrefix("/").Subrouter()
	publicRouter.NewRoute().
		Methods(http.MethodGet).
		Path("/public").
		HandlerFunc(handler.PublicPage)

	return r, nil
}
