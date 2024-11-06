package handlers

import (
	"context"
	"net/http"

	"github.com/burnerlee/compextAI/internal/logger"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

var s *Server

type Server struct {
	DB     *gorm.DB
	Ctx    context.Context
	Router *mux.Router
}

var err error

func InitServer(ctx context.Context) (*Server, error) {
	logger.GetLogger().Info("Initializing server")

	s = &Server{
		Ctx: ctx,
	}

	s.Router = mux.NewRouter()

	// add logger middleware to the router
	s.Router.Use(logger.LoggerMiddleware)

	// initialize the database
	logger.GetLogger().Info("Initializing database")
	s.DB, err = InitDB()
	if err != nil {
		logger.GetLogger().Errorf("Error initializing database: %v", err)
		return nil, err
	}

	logger.GetLogger().Info("Migrating database")
	if err := MigrateDB(s.DB); err != nil {
		logger.GetLogger().Errorf("Error migrating database: %v", err)
		return nil, err
	}

	logger.GetLogger().Info("Database initialized successfully")

	s.InitRoutes()

	return s, nil
}

func (s *Server) Run(addr string) {
	logger.GetLogger().Infof("Running server on %s", addr)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	handler := c.Handler(s.Router)
	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.GetLogger().Fatalf("Error starting server: %v", err)
	}

	logger.GetLogger().Info("Shutting down server gracefully")
}
