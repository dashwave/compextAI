package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/burnerlee/compextAI/controllers"
	"github.com/burnerlee/compextAI/models"
	"github.com/burnerlee/compextAI/utils/responses"
	"github.com/gorilla/mux"
)

func (s *Server) ListThreads(w http.ResponseWriter, r *http.Request) {
	// find all the threads from the db
	threads, err := models.GetAllThreads(s.DB)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, threads)
}

func (s *Server) CreateThread(w http.ResponseWriter, r *http.Request) {
	var request CreateThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	threadCreated, err := controllers.CreateThread(s.DB, &controllers.CreateThreadRequest{
		Title:    request.Title,
		Metadata: request.Metadata,
	})
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusCreated, threadCreated)
}

func (s *Server) GetThread(w http.ResponseWriter, r *http.Request) {
	threadID := mux.Vars(r)["id"]

	if threadID == "" {
		responses.Error(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	thread, err := models.GetThread(s.DB, threadID)
	if err != nil {
		responses.Error(w, http.StatusNotFound, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, thread)
}

func (s *Server) UpdateThread(w http.ResponseWriter, r *http.Request) {
	threadID := mux.Vars(r)["id"]

	if threadID == "" {
		responses.Error(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	var thread models.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	thread.Identifier = threadID

	updatedThread, err := models.UpdateThread(s.DB, &thread)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, updatedThread)
}

func (s *Server) DeleteThread(w http.ResponseWriter, r *http.Request) {
	threadID := mux.Vars(r)["id"]

	if threadID == "" {
		responses.Error(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	if err := models.DeleteThread(s.DB, threadID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusNoContent, "Thread deleted successfully")
}
