package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

// insertMockData inserts mock users and products into the database

type Config struct {
	PostgreSQL struct {
		Host           string `envconfig:"DB_HOST"`
		Port           int    `envconfig:"DB_PORT"`
		User           string `envconfig:"DB_USER"`
		Password       string `envconfig:"DB_PASSWORD"`
		Database       string `envconfig:"DB_NAME"`
		SharedDatabase string `envconfig:"SHARED_DB_NAME"`
	}

	SignatureRequest struct {
		SecretKey string `envconfig:"SIGNATURE_REQUEST_SECRET_KEY"`
	}
	Debug struct {
		IsProfiling bool `envconfig:"DEBUG_IS_PROFILING"`
	}
}

func LoadConfig() *Config {

	var config Config

	// Override with environment variables
	err := envconfig.Process("", &config)
	if err != nil {
		fmt.Printf("Error processing environment variables: %s\n", err)
		return nil
	}

	// Use config as needed
	fmt.Println("Signature Request Secret Key:", config.SignatureRequest.SecretKey)
	fmt.Println("Debug Is Profiling:", config.Debug.IsProfiling)

	return &config
}
