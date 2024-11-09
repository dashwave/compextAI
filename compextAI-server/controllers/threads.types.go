package controllers

type CreateThreadRequest struct {
	UserID    uint                   `json:"user_id"`
	ProjectID string                 `json:"project_id"`
	Title     string                 `json:"title"`
	Metadata  map[string]interface{} `json:"metadata"`
}
