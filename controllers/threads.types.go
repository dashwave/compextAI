package controllers

type CreateThreadRequest struct {
	Title    string            `json:"title"`
	Metadata map[string]string `json:"metadata"`
}
