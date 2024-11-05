package main

import (
	"context"
	"fmt"
	"os"

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

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	serverInstance.Run(fmt.Sprintf(":%s", serverPort))
}
