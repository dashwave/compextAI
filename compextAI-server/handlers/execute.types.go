package handlers

import "fmt"

type ExecuteThreadRequest struct {
	ExecutionModel          string      `json:"execution_model"`
	Temperature             float64     `json:"temperature"`
	Timeout                 int         `json:"timeout"`
	MaxCompletionTokens     int         `json:"max_completion_tokens"`
	TopP                    float64     `json:"top_p"`
	MaxOutputTokens         int         `json:"max_output_tokens"`
	ResponseFormat          interface{} `json:"response_format"`
	AppendAssistantResponse bool        `json:"append_assistant_response"`
}

func (r *ExecuteThreadRequest) Validate() error {
	if r.ExecutionModel == "" {
		return fmt.Errorf("execution_model is required")
	}
	return nil
}

type ExecuteThreadResponse struct {
	ThreadExecutionID string `json:"thread_execution_id"`
}

type ThreadExecutionStatusResponse struct {
	Status string `json:"status"`
}
