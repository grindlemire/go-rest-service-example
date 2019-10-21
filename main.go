package main

import (
	"io"
	"syscall"

	"github.com/grindlemire/log"
	"github.com/vrecan/death"

	"github.com/grindlemire/go-rest-service-example/pkg/config"
	"github.com/grindlemire/go-rest-service-example/pkg/rest"
	"github.com/grindlemire/go-rest-service-example/pkg/router"
)

func main() {
	log.Init(log.Default)
	opts, err := config.Load("./env")
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	router, err := router.NewRouter()
	if err != nil {
		log.Fatalf("unable to create path router: %v", err)
	}
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	goRoutines := []io.Closer{}

	// start the rest server for serving requests
	s := rest.NewServer(opts.ServePort, router)
	s.Start()
	goRoutines = append(goRoutines, s)

	err = d.WaitForDeath(goRoutines...)
	if err != nil {
		log.Fatalf("failed to shut down gracefully: %v", err)
	}
}
