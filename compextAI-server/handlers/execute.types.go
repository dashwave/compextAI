package handlers

import "fmt"

type ExecuteThreadRequest struct {
	ThreadExecutionParamsID string `json:"thread_execution_params_id"`
	AppendAssistantResponse bool   `json:"append_assistant_response"`
}

func (r *ExecuteThreadRequest) Validate() error {
	if r.ThreadExecutionParamsID == "" {
		return fmt.Errorf("thread_execution_params_id is required")
	}
	return nil
}

type ExecuteThreadResponse struct {
	ThreadExecutionID string `json:"thread_execution_id"`
}

type ThreadExecutionStatusResponse struct {
	Status string `json:"status"`
}
