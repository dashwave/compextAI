package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	"github.com/burnerlee/compextAI/models"
	"github.com/burnerlee/compextAI/utils"
	"github.com/burnerlee/compextAI/utils/responses"
	"github.com/gorilla/mux"
)

func (s *Server) ListThreadExecutionParams(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	projectName := mux.Vars(r)["projectname"]
	if projectName == "" {
		responses.Error(w, http.StatusBadRequest, "Project name is required")
		return
	}

	projectID, err := utils.GetProjectIDFromName(s.DB, projectName, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	executionParams, err := models.GetAllThreadExecutionParams(s.DB, uint(userID), projectID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make(ExecuteParamsResponse, 0)
	for _, executionParam := range executionParams {
		response = append(response, &squashedThreadExecutionParams{
			ProjectID:           executionParam.ProjectID,
			Identifier:          executionParam.Identifier,
			Name:                executionParam.Name,
			Environment:         executionParam.Environment,
			Model:               executionParam.Template.Model,
			Temperature:         executionParam.Template.Temperature,
			Timeout:             executionParam.Template.Timeout,
			MaxTokens:           executionParam.Template.MaxTokens,
			MaxCompletionTokens: executionParam.Template.MaxCompletionTokens,
			MaxOutputTokens:     executionParam.Template.MaxOutputTokens,
			TopP:                executionParam.Template.TopP,
			ResponseFormat:      executionParam.Template.ResponseFormat,
			SystemPrompt:        executionParam.Template.SystemPrompt,
		})
	}

	responses.JSON(w, http.StatusOK, response)
}

func (s *Server) CreateThreadExecutionParams(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request CreateThreadExecutionParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	projectID, err := utils.GetProjectIDFromName(s.DB, request.ProjectName, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking for existing execution params with the same name
	_, err = models.GetThreadExecutionParamsByUserIDAndNameAndEnvironment(s.DB, uint(userID), request.Name, request.Environment, projectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// no existing execution params with the same name
		} else {
			responses.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		responses.Error(w, http.StatusBadRequest, "Execution params with the same name already exists in this environment")
		return
	}

	executionParams := models.ThreadExecutionParams{
		UserID:      uint(userID),
		ProjectID:   projectID,
		Name:        request.Name,
		Environment: request.Environment,
		TemplateID:  request.TemplateID,
	}

	executionParamsCreated, err := models.CreateThreadExecutionParams(s.DB, &executionParams)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, executionParamsCreated)
}

func (s *Server) GetThreadExecutionParamsByNameAndEnv(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request GetThreadExecutionParamsByNameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	projectID, err := utils.GetProjectIDFromName(s.DB, request.ProjectName, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// no need to check access, because the user can only get his own execution params

	executionParams, err := models.GetThreadExecutionParamsByUserIDAndNameAndEnvironment(s.DB, uint(userID), request.Name, request.Environment, projectID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := &squashedThreadExecutionParams{
		Identifier:          executionParams.Identifier,
		Name:                executionParams.Name,
		Environment:         executionParams.Environment,
		Model:               executionParams.Template.Model,
		Temperature:         executionParams.Template.Temperature,
		Timeout:             executionParams.Template.Timeout,
		MaxTokens:           executionParams.Template.MaxTokens,
		MaxCompletionTokens: executionParams.Template.MaxCompletionTokens,
		MaxOutputTokens:     executionParams.Template.MaxOutputTokens,
		TopP:                executionParams.Template.TopP,
		ResponseFormat:      executionParams.Template.ResponseFormat,
		SystemPrompt:        executionParams.Template.SystemPrompt,
	}

	responses.JSON(w, http.StatusOK, response)
}

func (s *Server) DeleteThreadExecutionParams(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request DeleteThreadExecutionParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	projectID, err := utils.GetProjectIDFromName(s.DB, request.ProjectName, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	existingExecutionParams, err := models.GetThreadExecutionParamsByUserIDAndNameAndEnvironment(s.DB, uint(userID), request.Name, request.Environment, projectID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := models.DeleteThreadExecutionParams(s.DB, existingExecutionParams.Identifier); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, "Execution params deleted")
}

func (s *Server) ListThreadExecutionParamsTemplates(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	projectName := mux.Vars(r)["projectname"]
	if projectName == "" {
		responses.Error(w, http.StatusBadRequest, "Project name is required")
		return
	}

	threadExecutionParamsTemplates, err := models.GetAllThreadExecutionParamsTemplates(s.DB, uint(userID), projectName)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, threadExecutionParamsTemplates)
}

func (s *Server) CreateThreadExecutionParamsTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request CreateThreadExecutionParamsTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	responseFormat, err := json.Marshal(request.ResponseFormat)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	threadExecutionParamsTemplate := models.ThreadExecutionParamsTemplate{
		Name:                request.Name,
		UserID:              uint(userID),
		Model:               request.Model,
		Temperature:         request.Temperature,
		Timeout:             request.Timeout,
		MaxTokens:           request.MaxTokens,
		MaxCompletionTokens: request.MaxCompletionTokens,
		MaxOutputTokens:     request.MaxOutputTokens,
		SystemPrompt:        request.SystemPrompt,
		ResponseFormat:      responseFormat,
	}

	threadExecutionParamsTemplateCreated, err := models.CreateThreadExecutionParamsTemplate(s.DB, &threadExecutionParamsTemplate)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, threadExecutionParamsTemplateCreated)
}

func (s *Server) GetThreadExecutionParamsTemplateByID(w http.ResponseWriter, r *http.Request) {
	templateID := mux.Vars(r)["id"]
	if templateID == "" {
		responses.Error(w, http.StatusBadRequest, "Template ID is required")
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	hasAccess, err := utils.CheckThreadExecutionParamsTemplateAccess(s.DB, templateID, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !hasAccess {
		responses.Error(w, http.StatusForbidden, "You do not have access to this template")
		return
	}

	threadExecutionParamsTemplate, err := models.GetThreadExecutionParamsTemplateByID(s.DB, templateID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, threadExecutionParamsTemplate)
}

func (s *Server) DeleteThreadExecutionParamsTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := mux.Vars(r)["id"]
	if templateID == "" {
		responses.Error(w, http.StatusBadRequest, "Template ID is required")
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	hasAccess, err := utils.CheckThreadExecutionParamsTemplateAccess(s.DB, templateID, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !hasAccess {
		responses.Error(w, http.StatusForbidden, "You do not have access to this template")
		return
	}

	// check if there are any execution params using this template
	executionParams, err := models.GetThreadExecutionParamsByTemplateID(s.DB, templateID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			responses.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if len(executionParams) > 0 {
		responses.Error(w, http.StatusBadRequest, "Cannot delete template with execution params, execution params depend on this template")
		return
	}

	if err := models.DeleteThreadExecutionParamsTemplate(s.DB, templateID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, "Template deleted")
}

func (s *Server) UpdateThreadExecutionParamsTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := mux.Vars(r)["id"]
	if templateID == "" {
		responses.Error(w, http.StatusBadRequest, "Template ID is required")
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	hasAccess, err := utils.CheckThreadExecutionParamsTemplateAccess(s.DB, templateID, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !hasAccess {
		responses.Error(w, http.StatusForbidden, "You do not have access to this template")
		return
	}

	var request UpdateThreadExecutionParamsTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	threadExecutionParamsTemplate, err := models.GetThreadExecutionParamsTemplateByID(s.DB, templateID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseFormat, err := json.Marshal(request.ResponseFormat)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	threadExecutionParamsTemplate.Name = request.Name
	threadExecutionParamsTemplate.Model = request.Model
	threadExecutionParamsTemplate.Temperature = request.Temperature
	threadExecutionParamsTemplate.Timeout = request.Timeout
	threadExecutionParamsTemplate.MaxTokens = request.MaxTokens
	threadExecutionParamsTemplate.MaxCompletionTokens = request.MaxCompletionTokens
	threadExecutionParamsTemplate.MaxOutputTokens = request.MaxOutputTokens
	threadExecutionParamsTemplate.SystemPrompt = request.SystemPrompt
	threadExecutionParamsTemplate.ResponseFormat = responseFormat

	if err := models.UpdateThreadExecutionParamsTemplate(s.DB, threadExecutionParamsTemplate); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, "Template updated")
}
