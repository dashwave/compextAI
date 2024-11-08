package controllers

import "github.com/burnerlee/compextAI/models"

type ExecuteThreadRequest struct {
	UserID                         uint
	ThreadID                       string
	ThreadExecutionParamTemplateID string
	ThreadExecutionSystemPrompt    string
	AppendAssistantResponse        bool
	Messages                       []*models.Message
	FetchMessagesFromThread        bool
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
