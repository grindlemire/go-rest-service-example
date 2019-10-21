package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grindlemire/log"
	"github.com/rs/cors"
	"github.com/vrecan/life"
)

// Server is a wrapper around the http server that manages signals
// Note that I don't use life' lifecycle here because we have a blocking call for
// run (so I don't use life.Close or life.Done for managing the background thread.
// I use the server.ListenAndServe and server.Shutdown).
type Server struct {
	*life.Life
	server *http.Server
}

// NewServer creates a new http server with a router
func NewServer(port int, handler http.Handler) (s Server) {
	// This is where cors header stuff would be inserted
	c := cors.New(cors.Options{})
	s = Server{
		Life: life.NewLife(),
		server: &http.Server{
			Handler:      c.Handler(handler),
			Addr:         fmt.Sprintf(":%d", port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
	s.SetRun(s.run)
	return s
}

func (s Server) run() {
	log.Infof("server listening on [%s]", s.server.Addr)
	log.Infof("prometheus metrics at [%s/metrics]", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("unable to start listening: %v", err)
	}
}

// Close closes the server down gracefully
func (s Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
	log.Info("successfully shut down http server")
	return nil
}
