package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/grindlemire/log"

	"github.com/grindlemire/go-rest-service-example/pkg/config"
	"github.com/vrecan/life"
)

// Server is a wrapper around the http server that manages signals
type Server struct {
	*life.Life
	server *http.Server
}

// NewServer creates a new http server with a router
func NewServer(opts config.Opts, handler http.Handler) (s Server) {
	// This is where cors header stuff would be inserted
	c := cors.New(cors.Options{})

	s = Server{
		Life: life.NewLife(),
		server: &http.Server{
			Handler:      c.Handler(handler),
			Addr:         fmt.Sprintf("127.0.0.1:%d", opts.Port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
	s.SetRun(s.run)
	return s
}

func (s Server) run() {
	log.Infof("server starting to listen on %s", s.server.Addr)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("unable to start listening: %v", err)
		}
	}()

	for {
		select {
		case <-s.Done:
			return
		}
	}
}

// Close closes the server down gracefully
func (s Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
	s.Life.Close()
	log.Info("successfully shut down http server")
	return nil
}
