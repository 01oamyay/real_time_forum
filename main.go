package main

import (
	"log"

	"rlf/internal/app"
	"rlf/pkg/config"
)

// main loads the configuration and starts the HTTP server.
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
