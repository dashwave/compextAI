package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/burnerlee/compextAI/constants"
	"github.com/burnerlee/compextAI/controllers"
	"github.com/burnerlee/compextAI/models"
	"github.com/burnerlee/compextAI/utils"
	"github.com/burnerlee/compextAI/utils/responses"
	"github.com/gorilla/mux"
)

func (s *Server) ExecuteThread(w http.ResponseWriter, r *http.Request) {
	threadID := mux.Vars(r)["id"]

	if threadID == "" {
		responses.Error(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	var request ExecuteThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(threadID); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	if threadID != constants.THREAD_IDENTIFIER_FOR_NULL_THREAD {
		hasAccess, err := utils.CheckThreadAccess(s.DB, threadID, uint(userID))
		if err != nil {
			responses.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !hasAccess {
			responses.Error(w, http.StatusForbidden, "You are not authorized to execute this thread")
			return
		}
	}

	threadExecution, err := controllers.ExecuteThread(s.DB, &controllers.ExecuteThreadRequest{
		UserID:                      uint(userID),
		ThreadID:                    threadID,
		ThreadExecutionParamID:      request.ThreadExecutionParamID,
		AppendAssistantResponse:     request.AppendAssistantResponse,
		ThreadExecutionSystemPrompt: request.ThreadExecutionSystemPrompt,
		Messages:                    request.Messages,
	})
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, threadExecution)
}

func (s *Server) GetThreadExecutionStatus(w http.ResponseWriter, r *http.Request) {
	executionID := mux.Vars(r)["id"]

	if executionID == "" {
		responses.Error(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	hasAccess, err := utils.CheckThreadExecutionAccess(s.DB, executionID, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !hasAccess {
		responses.Error(w, http.StatusForbidden, "You are not authorized to access this thread execution")
		return
	}

	threadExecution, err := models.GetThreadExecutionByID(s.DB, executionID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, ThreadExecutionStatusResponse{Status: threadExecution.Status})
}

func (s *Server) GetThreadExecutionResponse(w http.ResponseWriter, r *http.Request) {
	executionID := mux.Vars(r)["id"]

	if executionID == "" {
		responses.Error(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	userID, err := utils.GetUserIDFromRequest(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	hasAccess, err := utils.CheckThreadExecutionAccess(s.DB, executionID, uint(userID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !hasAccess {
		responses.Error(w, http.StatusForbidden, "You are not authorized to access this thread execution")
		return
	}

	threadExecution, err := models.GetThreadExecutionByID(s.DB, executionID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if threadExecution.Status != models.ThreadExecutionStatus_COMPLETED {
		responses.Error(w, http.StatusBadRequest, fmt.Sprintf("Thread execution is: %s", threadExecution.Status))
		return
	}

	var responseContent interface{}
	if err := json.Unmarshal(threadExecution.Output, &responseContent); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"response": responseContent,
		"content":  threadExecution.Content,
		"role":     threadExecution.Role,
	})
}
