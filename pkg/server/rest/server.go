package rest

import (
	"context"
	"fmt"
	"io/ioutil"
	stdlogger "log"
	"net/http"
	"strings"
	"time"

	"github.com/grindlemire/log"
	"github.com/vrecan/life"
)

// This is really dumb but it keeps the http server from logging that we have an unknown certificate
// using the std logger. This could/probably should be redirected to my logger
func init() {
	stdlogger.SetOutput(ioutil.Discard)
}

// Server is a wrapper around the http server that manages signals
// Note that I don't use life' lifecycle here because we have a blocking call for
// run (so I don't use life.Close or life.Done for managing the background thread.
// I use the server.ListenAndServe and server.Shutdown).
type Server struct {
	*life.Life
	httpsServer        *http.Server
	httpRedirectServer *http.Server
	tlsCertPath        string
	tlsKeyPath         string
}

// NewServer creates a new http server with a router. The reason why we pass through is because we are using
// a functional constructor and we don't want to pollute the main struct with intermediate config state
func NewServer(opts ...Opt) (s *Server, err error) {
	return build(opts...)
}

func (s Server) run() {
	go func() {
		log.Infof("http redirect server listening on [%s]", s.httpRedirectServer.Addr)
		err := s.httpRedirectServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("unable to start listening on http redirect: %v", err)
		}
	}()

	log.Infof("server listening on [%s]", s.httpsServer.Addr)
	log.Infof("prometheus metrics at [%s/metrics]", s.httpsServer.Addr)

	err := s.httpsServer.ListenAndServeTLS(s.tlsCertPath, s.tlsKeyPath)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("unable to start listening: %v", err)
	}
}

func (s Server) getHTTPSServerAddr() string {
	return s.httpsServer.Addr
}

// Close closes the server down gracefully
func (s Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.httpsServer.Shutdown(ctx)
	s.httpRedirectServer.Shutdown(ctx)
	log.Info("successfully shut down http server")
	return nil
}

// createHTTPSRedirect creates a redirect function that will redirect us to the right rest server
// to redirect us from http to https, even if the https server is served on a nonstandard port
func createHTTPSRedirect(httpsPort int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cleanedHost := strings.Split(r.Host, ":")[0]
		from := fmt.Sprintf("http://%s%s", r.Host, r.RequestURI)
		redirect := fmt.Sprintf("https://%s:%d%s", cleanedHost, httpsPort, r.RequestURI)
		log.Infof("Redirecting [%s] to [%s]", from, redirect)
		http.Redirect(w, r, redirect, http.StatusMovedPermanently)
	}
}
