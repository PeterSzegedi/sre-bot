package env

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

// Populate tokens from environment variables
func PopulateEnv() EnvConfig {
	var env EnvConfig
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("Failed to read env vars: %s", err)
	}
	return env
}
