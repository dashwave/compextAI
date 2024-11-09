package handlers

import "errors"

type CreateThreadRequest struct {
	Title       string                 `json:"title"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProjectName string                 `json:"project_name"`
}

func (r *CreateThreadRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	if r.ProjectName == "" {
		return errors.New("project name is required")
	}
	return nil
}

type UpdateThreadRequest struct {
	Title    string                 `json:"title"`
	Metadata map[string]interface{} `json:"metadata"`
}
