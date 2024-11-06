package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/burnerlee/compextAI/internal/logger"
	"github.com/burnerlee/compextAI/internal/providers/chat"
	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

func ExecuteThread(db *gorm.DB, req *ExecuteThreadRequest) (interface{}, error) {
	chatProvider, err := chat.GetChatCompletionsProvider(req.ExecutionModel)
	if err != nil {
		logger.GetLogger().Errorf("Error getting chat provider: %s: %v", req.ExecutionModel, err)
		return nil, err
	}

	// get the thread
	thread, err := models.GetThread(db, req.ThreadID)
	if err != nil {
		logger.GetLogger().Errorf("Error getting thread: %s: %v", req.ThreadID, err)
		return nil, err
	}

	threadExecutionParams := &models.ThreadExecutionParams{
		Model:               req.ExecutionModel,
		Temperature:         req.Temperature,
		Timeout:             req.Timeout,
		MaxCompletionTokens: req.MaxCompletionTokens,
		TopP:                req.TopP,
		MaxOutputTokens:     req.MaxOutputTokens,
		ResponseFormat:      req.ResponseFormat,
	}

	// start thread execution
	threadExecutionParamsBytes, err := json.Marshal(threadExecutionParams)
	if err != nil {
		logger.GetLogger().Errorf("Error marshalling thread execution params: %v", err)
		return nil, err
	}

	threadExecution := models.ThreadExecution{
		ThreadID:              req.ThreadID,
		ThreadExecutionParams: threadExecutionParamsBytes,
		Status:                models.ThreadExecutionStatus_IN_PROGRESS,
	}

	if err := models.CreateThreadExecution(db, &threadExecution); err != nil {
		logger.GetLogger().Errorf("Error creating thread execution: %v", err)
		return nil, err
	}

	go func(thread *models.Thread, threadExecutionParams *models.ThreadExecutionParams, appendAssistantResponse bool) {
		// get the user
		user, err := models.GetUserByID(db, thread.UserID)
		if err != nil {
			logger.GetLogger().Errorf("Error getting user: %d: %v", thread.UserID, err)
			return
		}

		// execute the thread using the chat provider
		statusCode, threadExecutionResponse, err := chatProvider.ExecuteThread(db, user.OpenAIKey, thread, threadExecutionParams)
		if err != nil {
			logger.GetLogger().Errorf("Error executing thread: %s: %v: %v", req.ThreadID, err, threadExecutionResponse)
			handleThreadExecutionError(db, &threadExecution, fmt.Errorf("error executing thread: %v: %v", err, threadExecutionResponse))
			return
		}

		if statusCode != http.StatusOK {
			logger.GetLogger().Errorf("Error executing thread: %s: status code: %d: %v", req.ThreadID, statusCode, threadExecutionResponse)
			handleThreadExecutionError(db, &threadExecution, fmt.Errorf("status code: %d: %v", statusCode, threadExecutionResponse))
			return
		}

		logger.GetLogger().Infof("Thread execution completed: %s", req.ThreadID)
		handleThreadExecutionSuccess(db, &threadExecution, threadExecutionResponse, appendAssistantResponse)
	}(thread, threadExecutionParams, req.AppendAssistantResponse)

	return threadExecution, nil
}

func handleThreadExecutionError(db *gorm.DB, threadExecution *models.ThreadExecution, err error) {
	threadExecution.Status = models.ThreadExecutionStatus_FAILED
	threadExecution.Output = err.Error()
	models.UpdateThreadExecution(db, threadExecution)
}

func handleThreadExecutionSuccess(db *gorm.DB, threadExecution *models.ThreadExecution, threadExecutionResponse interface{}, appendAssistantResponse bool) {
	threadExecution.Status = models.ThreadExecutionStatus_COMPLETED
	threadExecution.Output = fmt.Sprintf("%v", threadExecutionResponse)
	models.UpdateThreadExecution(db, threadExecution)

	if appendAssistantResponse {
		responseMessage := threadExecutionResponse.(map[string]interface{})

		if err := models.CreateMessage(db, &models.Message{
			ThreadID: threadExecution.ThreadID,
			Role:     responseMessage["role"].(string),
			Content:  responseMessage["content"].(string),
		}); err != nil {
			logger.GetLogger().Errorf("Error creating assistant message: %v", err)
		} else {
			logger.GetLogger().Infof("assistant message created")
		}
	}
}
