package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/burnerlee/compextAI/models"
	"github.com/burnerlee/compextAI/utils"
	"github.com/burnerlee/compextAI/utils/responses"
	"github.com/gorilla/mux"
)

func (s *Server) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// check if the project name is already taken
	if _, err := models.GetProjectByName(s.DB, request.Name, uint(userID)); err == nil {
		responses.Error(w, http.StatusBadRequest, "project name already taken")
		return
	}

	project := &models.Project{
		UserID:      uint(userID),
		Name:        request.Name,
		Description: request.Description,
	}

	if err := models.CreateProject(s.DB, project); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, project)
}

func (s *Server) GetProject(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["id"]

	if projectID == "" {
		responses.Error(w, http.StatusBadRequest, "project id is required")
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	hasAccess, err := utils.CheckProjectAccess(s.DB, projectID, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !hasAccess {
		responses.Error(w, http.StatusForbidden, "you do not have access to this project")
		return
	}

	project, err := models.GetProject(s.DB, projectID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, project)
}

func (s *Server) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	projects, err := models.GetAllProjects(s.DB, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, projects)
}

func (s *Server) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["id"]

	if projectID == "" {
		responses.Error(w, http.StatusBadRequest, "project id is required")
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	hasAccess, err := utils.CheckProjectAccess(s.DB, projectID, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !hasAccess {
		responses.Error(w, http.StatusForbidden, "you do not have access to this project")
		return
	}

	if err := models.DeleteProject(s.DB, projectID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, "Project deleted successfully")
}

func (s *Server) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["id"]

	if projectID == "" {
		responses.Error(w, http.StatusBadRequest, "project id is required")
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	hasAccess, err := utils.CheckProjectAccess(s.DB, projectID, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !hasAccess {
		responses.Error(w, http.StatusForbidden, "you do not have access to this project")
		return
	}

	var request UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	existingProject, err := models.GetProject(s.DB, projectID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := models.UpdateProject(s.DB, &models.Project{
		Base: models.Base{
			ID: existingProject.ID,
		},
		Name:        request.Name,
		Description: request.Description,
	}); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, existingProject)
}
