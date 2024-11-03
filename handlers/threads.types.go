package handlers

import "errors"

type CreateThreadRequest struct {
	Title    string            `json:"title"`
	Metadata map[string]string `json:"metadata"`
}

func (r *CreateThreadRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	return nil
}

type UpdateThreadRequest struct {
	Title    string            `json:"title"`
	Metadata map[string]string `json:"metadata"`
}
