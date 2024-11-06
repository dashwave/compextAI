package controllers

type ExecuteThreadRequest struct {
	ThreadID                string
	ThreadExecutionParamsID string
	AppendAssistantResponse bool
}

type ExecuteThreadResponse struct {
	ThreadExecutionID string
}
