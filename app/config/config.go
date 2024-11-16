package config

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Println("PORT environment variable not set, using default port 8000")
	}

	return &Config{
		Port: port,
	}
}
