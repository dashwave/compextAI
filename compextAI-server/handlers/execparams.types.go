package handlers

import "errors"

type CreateThreadExecutionParamsRequest struct {
	Name                string      `json:"name"`
	Environment         string      `json:"environment"`
	Model               string      `json:"model"`
	Temperature         float64     `json:"temperature"`
	Timeout             int         `json:"timeout"`
	MaxTokens           int         `json:"max_tokens"`
	MaxCompletionTokens int         `json:"max_completion_tokens"`
	TopP                float64     `json:"top_p"`
	MaxOutputTokens     int         `json:"max_output_tokens"`
	ResponseFormat      interface{} `json:"response_format"`
	SystemPrompt        string      `json:"system_prompt"`
}

func (r *CreateThreadExecutionParamsRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Environment == "" {
		return errors.New("environment is required")
	}
	if r.Model == "" {
		return errors.New("model is required")
	}
	return nil
}

type UpdateThreadExecutionParamsRequest struct {
	CreateThreadExecutionParamsRequest
}

func (r *UpdateThreadExecutionParamsRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Environment == "" {
		return errors.New("environment is required")
	}
	return nil
}

type GetThreadExecutionParamsByNameRequest struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
}

func (r *GetThreadExecutionParamsByNameRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Environment == "" {
		return errors.New("environment is required")
	}
	return nil
}

type DeleteThreadExecutionParamsRequest struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
}

func (r *DeleteThreadExecutionParamsRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Environment == "" {
		return errors.New("environment is required")
	}
	return nil
}
