package main

import (
	"io"
	"syscall"

	"github.com/grindlemire/log"
	"github.com/vrecan/death"

	"github.com/grindlemire/go-rest-service-example/pkg/server/config"
	"github.com/grindlemire/go-rest-service-example/pkg/server/endpoint"
	"github.com/grindlemire/go-rest-service-example/pkg/server/rest"
	"github.com/grindlemire/go-rest-service-example/pkg/server/router"
)

func main() {
	log.Init(log.Opts{})
	opts, err := config.Load("./env")
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	router, err := router.NewRouter(endpoint.CreateEndpoints())
	if err != nil {
		log.Fatalf("unable to create path router: %v", err)
	}
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	goRoutines := []io.Closer{}

	// start the rest server for serving requests
	s, err := rest.NewServer(
		rest.Handler(router),
		rest.HTTPPort(opts.HTTPPort),
		rest.HTTPSPort(opts.HTTPSPort),
		rest.TLSCertPath(opts.TLSCertPath),
		rest.TLSKeyPath(opts.TLSKeyPath),
	)
	if err != nil {
		log.Fatalf("unable to create rest server: %v", err)
	}
	s.Start()
	goRoutines = append(goRoutines, s)

	err = d.WaitForDeath(goRoutines...)
	if err != nil {
		log.Fatalf("failed to shut down gracefully: %v", err)
	}
}
