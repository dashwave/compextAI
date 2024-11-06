package controllers

type ExecuteThreadRequest struct {
	ThreadID                string
	ExecutionModel          string
	Temperature             float64
	Timeout                 int
	MaxCompletionTokens     int
	TopP                    float64
	MaxOutputTokens         int
	ResponseFormat          interface{}
	AppendAssistantResponse bool
	SystemPrompt            string
}

type ExecuteThreadResponse struct {
	ThreadExecutionID string
}
