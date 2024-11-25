package controllers

import (
	"encoding/json"

	"github.com/burnerlee/compextAI/models"
)

type ExecuteThreadRequest struct {
	UserID                         uint
	ThreadID                       string
	ThreadExecutionParamTemplateID string
	ThreadExecutionSystemPrompt    string
	AppendAssistantResponse        bool
	Messages                       []*models.Message
	FetchMessagesFromThread        bool
	ProjectID                      string
	Metadata                       json.RawMessage
}

type ExecuteThreadResponse struct {
	ThreadExecutionID string
}

type RerunThreadExecutionRequest struct {
	UserID                         uint
	ExecutionID                    string
	ThreadExecutionParamTemplateID string
	SystemPrompt                   string
	AppendAssistantResponse        bool
}
