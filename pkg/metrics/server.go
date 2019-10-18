package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grindlemire/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vrecan/life"
)

// Server is the metrics prometheus server
type Server struct {
	*life.Life

	server *http.Server
}

// NewServer creates a new prometheus metrics server
func NewServer(port int) (s *Server) {
	r := http.NewServeMux()
	r.Handle("/metrics", promhttp.Handler())

	s = &Server{
		Life: life.NewLife(),
		server: &http.Server{
			Handler:      r,
			Addr:         fmt.Sprintf("127.0.0.1:%d", port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
	s.SetRun(s.run)
	return s
}

func (s Server) run() {
	go func() {
		log.Infof("Metrics server running on [%v]", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start metrics server: %v", err)
		}
		log.Info("Metrics server has shut down")
	}()

	for {
		select {
		case <-s.Done:
			return
		}
	}
}

// Close closes down the metrics server
func (s Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
	s.Life.Close()
	return nil
}
