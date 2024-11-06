package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	"github.com/burnerlee/compextAI/models"
	"github.com/burnerlee/compextAI/utils"
	"github.com/burnerlee/compextAI/utils/responses"
)

func (s *Server) ListThreadExecutionParams(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	executionParams, err := models.GetAllThreadExecutionParams(s.DB, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, executionParams)
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

	// checking for existing execution params with the same name
	_, err = models.GetThreadExecutionParamsByUserIDAndNameAndEnvironment(s.DB, uint(userID), request.Name, request.Environment)
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

	responseFormat, err := json.Marshal(request.ResponseFormat)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	executionParams := models.ThreadExecutionParams{
		UserID:              uint(userID),
		Name:                request.Name,
		Environment:         request.Environment,
		Model:               request.Model,
		Temperature:         request.Temperature,
		Timeout:             request.Timeout,
		MaxTokens:           request.MaxTokens,
		MaxCompletionTokens: request.MaxCompletionTokens,
		MaxOutputTokens:     request.MaxOutputTokens,
		SystemPrompt:        request.SystemPrompt,
		ResponseFormat:      responseFormat,
		TopP:                request.TopP,
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

	// no need to check access, because the user can only get his own execution params

	executionParams, err := models.GetThreadExecutionParamsByUserIDAndNameAndEnvironment(s.DB, uint(userID), request.Name, request.Environment)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, executionParams)
}

func (s *Server) UpdateThreadExecutionParams(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var request UpdateThreadExecutionParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	existingExecutionParams, err := models.GetThreadExecutionParamsByUserIDAndNameAndEnvironment(s.DB, uint(userID), request.Name, request.Environment)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseFormatJson, err := json.Marshal(request.ResponseFormat)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	existingExecutionParams.Name = request.Name
	existingExecutionParams.Model = request.Model
	existingExecutionParams.Temperature = request.Temperature
	existingExecutionParams.Timeout = request.Timeout
	existingExecutionParams.MaxTokens = request.MaxTokens
	existingExecutionParams.MaxCompletionTokens = request.MaxCompletionTokens
	existingExecutionParams.MaxOutputTokens = request.MaxOutputTokens
	existingExecutionParams.SystemPrompt = request.SystemPrompt
	existingExecutionParams.ResponseFormat = responseFormatJson
	existingExecutionParams.TopP = request.TopP

	if err := models.UpdateThreadExecutionParams(s.DB, existingExecutionParams); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, "Execution params updated")
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

	existingExecutionParams, err := models.GetThreadExecutionParamsByUserIDAndNameAndEnvironment(s.DB, uint(userID), request.Name, request.Environment)
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
