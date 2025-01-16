package main

import (
	"log"

	"rlf/internal/app"
	"rlf/pkg/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
