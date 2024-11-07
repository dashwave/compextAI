package handlers

import "errors"

type CreateMessageRequest struct {
	Content  string                 `json:"content"`
	Role     string                 `json:"role"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (r *CreateMessageRequest) Validate() error {
	if r.Content == "" {
		return errors.New("content is required")
	}
	if r.Role == "" {
		return errors.New("role is required")
	}
	return nil
}

type UpdateMessageRequest struct {
	Content  string                 `json:"content"`
	Role     string                 `json:"role"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (r *UpdateMessageRequest) Validate() error {
	// TODO: validate the request
	return nil
}
