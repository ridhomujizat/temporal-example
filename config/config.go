package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// GetEnv TODO: Add validation for the config
func GetEnv() (config *Config, er error) {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading .env file")
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
