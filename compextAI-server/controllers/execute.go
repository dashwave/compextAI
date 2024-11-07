package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/burnerlee/compextAI/constants"
	"github.com/burnerlee/compextAI/internal/logger"
	"github.com/burnerlee/compextAI/internal/providers/chat"
	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

func ExecuteThread(db *gorm.DB, req *ExecuteThreadRequest) (interface{}, error) {
	threadExecutionParams, err := models.GetThreadExecutionParamsByID(db, req.ThreadExecutionParamID)
	if err != nil {
		logger.GetLogger().Errorf("Error getting thread execution params: %s: %v", req.ThreadExecutionParamID, err)
		return nil, err
	}

	if req.ThreadExecutionSystemPrompt != "" {
		logger.GetLogger().Infof("Setting thread execution system prompt: %s", req.ThreadExecutionSystemPrompt)
		threadExecutionParams.Template.SystemPrompt = req.ThreadExecutionSystemPrompt
	}

	chatProvider, err := chat.GetChatCompletionsProvider(threadExecutionParams.Template.Model)
	if err != nil {
		logger.GetLogger().Errorf("Error getting chat provider: %s: %v", threadExecutionParams.Template.Model, err)
		return nil, err
	}

	var messages []*models.Message
	if req.ThreadID != constants.THREAD_IDENTIFIER_FOR_NULL_THREAD {
		// get the thread
		threadMessages, err := models.GetAllMessages(db, req.ThreadID)
		if err != nil {
			logger.GetLogger().Errorf("Error getting thread: %s: %v", req.ThreadID, err)
			return nil, err
		}
		messages = threadMessages
	} else {
		messages = req.Messages
	}

	threadExecution := &models.ThreadExecution{
		UserID:                 req.UserID,
		ThreadID:               req.ThreadID,
		ThreadExecutionParamID: req.ThreadExecutionParamID,
		Status:                 models.ThreadExecutionStatus_IN_PROGRESS,
	}

	threadExecution, err = models.CreateThreadExecution(db, threadExecution)
	if err != nil {
		logger.GetLogger().Errorf("Error creating thread execution: %v", err)
		return nil, err
	}

	go func(p chat.ChatCompletionsProvider, messages []*models.Message, threadExecution models.ThreadExecution, threadExecutionParams models.ThreadExecutionParams, appendAssistantResponse bool) {
		// get the user
		user, err := models.GetUserByID(db, threadExecution.UserID)
		if err != nil {
			logger.GetLogger().Errorf("Error getting user: %d: %v", threadExecution.UserID, err)
			return
		}

		// execute the thread using the chat provider
		statusCode, threadExecutionResponse, err := chatProvider.ExecuteThread(db, user, messages, &threadExecutionParams, threadExecution.Identifier)
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
		handleThreadExecutionSuccess(db, p, &threadExecution, threadExecutionResponse, appendAssistantResponse)
	}(chatProvider, messages, *threadExecution, *threadExecutionParams, req.AppendAssistantResponse)

	return threadExecution, nil
}

func handleThreadExecutionError(db *gorm.DB, threadExecution *models.ThreadExecution, err error) {
	threadExecution.Status = models.ThreadExecutionStatus_FAILED
	threadExecution.Output = json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	models.UpdateThreadExecution(db, threadExecution)
}

func handleThreadExecutionSuccess(db *gorm.DB, p chat.ChatCompletionsProvider, threadExecution *models.ThreadExecution, threadExecutionResponse interface{}, appendAssistantResponse bool) {
	responseJson, err := json.Marshal(threadExecutionResponse)
	if err != nil {
		logger.GetLogger().Errorf("Error marshalling thread execution response: %v", err)
		return
	}

	message, err := p.ConvertExecutionResponseToMessage(threadExecutionResponse)
	if err != nil {
		logger.GetLogger().Errorf("Error converting thread execution response to message: %v", err)
		return
	}

	if appendAssistantResponse {
		logger.GetLogger().Infof("Appending assistant response")

		if err := models.CreateMessage(db, &models.Message{
			ThreadID: threadExecution.ThreadID,
			Role:     message.Role,
			Content:  message.Content,
			Metadata: message.Metadata,
		}); err != nil {
			logger.GetLogger().Errorf("Error creating assistant message: %v", err)
		} else {
			logger.GetLogger().Infof("assistant message created")
		}
	}

	threadExecution.Status = models.ThreadExecutionStatus_COMPLETED
	threadExecution.Output = responseJson
	threadExecution.Content = message.Content
	threadExecution.Role = message.Role
	threadExecution.ExecutionResponseMetadata = message.Metadata
	models.UpdateThreadExecution(db, threadExecution)
}
