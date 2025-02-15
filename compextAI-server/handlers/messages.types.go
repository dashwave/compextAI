package handlers

import (
	"encoding/json"
	"errors"
	"time"
)

type messageResponse struct {
	Identifier   string          `json:"identifier"`
	Content      interface{}     `json:"content"`
	Role         string          `json:"role"`
	ToolCallID   string          `json:"tool_call_id"`
	ThreadID     string          `json:"thread_id"`
	Metadata     json.RawMessage `json:"metadata"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	ToolCalls    interface{}     `json:"tool_calls"`
	FunctionCall interface{}     `json:"function_call"`
}

type createMessage struct {
	Content      interface{}            `json:"content"`
	Role         string                 `json:"role"`
	ToolCallID   string                 `json:"tool_call_id"`
	Metadata     map[string]interface{} `json:"metadata"`
	ToolCalls    interface{}            `json:"tool_calls"`
	FunctionCall interface{}            `json:"function_call"`
}

func (m *createMessage) Validate() error {
	// check if interface is nil
	if m.Content == nil {
		return errors.New("content is required")
	}
	if m.Role == "" {
		return errors.New("role is required")
	}
	return nil
}

type CreateMessageRequest struct {
	Messages []*createMessage `json:"messages"`
}

func (r *CreateMessageRequest) Validate() error {
	if len(r.Messages) == 0 {
		return errors.New("at least one message is required")
	}
	for _, message := range r.Messages {
		if err := message.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type UpdateMessageRequest struct {
	Content      interface{}            `json:"content"`
	Role         string                 `json:"role"`
	ToolCallID   string                 `json:"tool_call_id"`
	Metadata     map[string]interface{} `json:"metadata"`
	ToolCalls    interface{}            `json:"tool_calls"`
	FunctionCall interface{}            `json:"function_call"`
}

func (r *UpdateMessageRequest) Validate() error {
	// TODO: validate the request
	return nil
}
