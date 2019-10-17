package main

import (
	"io"
	"syscall"

	"github.com/grindlemire/go-rest-service-example/pkg/server"
	"github.com/vrecan/death"

	"github.com/grindlemire/go-rest-service-example/pkg/router"
	"github.com/grindlemire/log"

	"github.com/grindlemire/go-rest-service-example/pkg/config"
)

func main() {
	log.Init(log.Default)
	opts, err := config.Load("./env")
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	router, err := router.NewRouter(opts)
	if err != nil {
		log.Fatalf("unable to create path router: %v", err)
	}
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	goRoutines := []io.Closer{}

	s := server.NewServer(opts, router)
	s.Start()
	goRoutines = append(goRoutines, s)

	err = d.WaitForDeath(goRoutines...)
	if err != nil {
		log.Fatalf("failed to shut down gracefully: %v", err)
	}
	log.Info("cleanly shut down")
}
