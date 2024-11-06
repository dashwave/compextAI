package handlers

import (
	"net/http"

	"github.com/burnerlee/compextAI/utils/responses"
)

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "pong")
}
