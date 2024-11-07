package controllers

type ExecuteThreadRequest struct {
	ThreadID                    string
	ThreadExecutionParamID      string
	ThreadExecutionSystemPrompt string
	AppendAssistantResponse     bool
}

type ExecuteThreadResponse struct {
	ThreadExecutionID string
}
