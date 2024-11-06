package handlers

import "fmt"

type ExecuteThreadRequest struct {
	ThreadExecutionParamID  string `json:"thread_execution_param_id"`
	AppendAssistantResponse bool   `json:"append_assistant_response"`
}

func (r *ExecuteThreadRequest) Validate() error {
	if r.ThreadExecutionParamID == "" {
		return fmt.Errorf("thread_execution_param_id is required")
	}
	return nil
}

type ExecuteThreadResponse struct {
	ThreadExecutionID string `json:"thread_execution_id"`
}

type ThreadExecutionStatusResponse struct {
	Status string `json:"status"`
}
