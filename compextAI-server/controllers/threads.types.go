package controllers

type CreateThreadRequest struct {
	UserID   uint              `json:"user_id"`
	Title    string            `json:"title"`
	Metadata map[string]string `json:"metadata"`
}
