package main

import (
	"context"
	"fmt"
	"os"

	"github.com/burnerlee/compextAI/handlers"
	"github.com/burnerlee/compextAI/internal/logger"
)

func main() {
	ctx := context.Background()

	serverInstance, err := handlers.InitServer(ctx)
	if err != nil {
		logger.GetLogger().Fatalf("Error initializing server: %v", err)
	}

	logger.GetLogger().Info("Server initialized successfully")

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	serverInstance.Run(fmt.Sprintf(":%s", serverPort))
}
