package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Opts ...
type Opts struct {
	Port int `required:"true"`
}

// Load the configuration
func Load(envfile string) (opts Opts, err error) {
	// load in the env file to the environment
	err = godotenv.Load(envfile)
	if err != nil {
		return opts, err
	}

	// parse our environment into the opts struct
	err = envconfig.Process("", &opts)
	if err != nil {
		return opts, err
	}

	return opts, err
}
