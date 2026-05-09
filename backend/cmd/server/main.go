package main

import (
	"flag"
	"log"

	"admin-platform/backend/internal/bootstrap"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "config file path")
	flag.Parse()

	app, err := bootstrap.New(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
