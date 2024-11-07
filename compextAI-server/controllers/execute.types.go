package controllers

import "github.com/burnerlee/compextAI/models"

type ExecuteThreadRequest struct {
	UserID                      uint
	ThreadID                    string
	ThreadExecutionParamID      string
	ThreadExecutionSystemPrompt string
	AppendAssistantResponse     bool
	Messages                    []*models.Message
}

type ExecuteThreadResponse struct {
	ThreadExecutionID string
}

type RerunThreadExecutionRequest struct {
	UserID                  uint
	ExecutionID             string
	ThreadExecutionParamID  string
	SystemPrompt            string
	AppendAssistantResponse bool
}
