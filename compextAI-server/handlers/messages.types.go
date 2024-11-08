package handlers

import "errors"

type createMessage struct {
	Content  string                 `json:"content"`
	Role     string                 `json:"role"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (m *createMessage) Validate() error {
	if m.Content == "" {
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
	Content  string                 `json:"content"`
	Role     string                 `json:"role"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (r *UpdateMessageRequest) Validate() error {
	// TODO: validate the request
	return nil
}
