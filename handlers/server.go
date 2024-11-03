package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
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
	log.Info("Initializing server")

	s = &Server{
		Ctx: ctx,
	}

	s.Router = mux.NewRouter()

	// initialize the database
	log.Info("Initializing database")
	s.DB, err = InitDB()
	if err != nil {
		log.Errorf("Error initializing database: %v", err)
		return nil, err
	}

	log.Info("Migrating database")
	if err := MigrateDB(s.DB); err != nil {
		log.Errorf("Error migrating database: %v", err)
		return nil, err
	}

	log.Info("Database initialized successfully")

	s.InitRoutes()

	return s, nil
}

func (s *Server) Run(addr string) {
	log.Infof("Running server on %s", addr)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(s.Router)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Info("Shutting down server gracefully")
}
