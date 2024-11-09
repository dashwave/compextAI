package handlers

import (
	"errors"
	"regexp"
	"strings"
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreateProjectRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	// name should not contain spaces or special characters
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", r.Name)
	if err != nil {
		return errors.New("invalid name")
	}
	if strings.ContainsAny(r.Name, " ") || !matched {
		return errors.New("name should not contain spaces or special characters")
	}
	return nil
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *UpdateProjectRequest) Validate() error {
	return nil
}
