package controllers

type CreateMessageRequest struct {
	ThreadID string            `json:"thread_id"`
	Content  string            `json:"content"`
	Role     string            `json:"role"`
	Metadata map[string]string `json:"metadata"`
}
