package main

import (
	"context"

	"github.com/burnerlee/compextAI/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	serverInstance, err := handlers.InitServer(ctx)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	log.Info("Server initialized successfully")
	serverInstance.Run(":8080")
}
