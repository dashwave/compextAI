package openai

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/burnerlee/compextAI/models"
)

const (
	GPT4O_MODEL      = "gpt-4o"
	GPT4O_OWNER      = "openai"
	GPT4O_IDENTIFIER = "gpt-4o"
)

type GPT4O struct {
	owner        string
	model        string
	allowedRoles []string
}

func NewGPT4O() *GPT4O {
	return &GPT4O{
		owner:        GPT4O_OWNER,
		model:        GPT4O_MODEL,
		allowedRoles: []string{"user", "assistant", "system"},
	}
}

func (g *GPT4O) GetProviderOwner() string {
	return g.owner
}

func (g *GPT4O) GetProviderModel() string {
	return g.model
}

func (g *GPT4O) GetProviderIdentifier() string {
	return GPT4O_IDENTIFIER
}

func (g *GPT4O) ValidateMessage(message *models.Message) error {
	if message.Content == "" {
		return fmt.Errorf("message content is empty")
	}

	if !slices.Contains(g.allowedRoles, message.Role) {
		return fmt.Errorf("message role is invalid, only %v are allowed", g.allowedRoles)
	}
	return nil
}

type gpt4oOpenAIMessage struct {
	Role     string            `json:"role"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

func (g *GPT4O) ConvertMessageToProviderFormat(message *models.Message) (interface{}, error) {
	var metadata map[string]string
	if err := json.Unmarshal(message.Metadata, &metadata); err != nil {
		return nil, err
	}

	return gpt4oOpenAIMessage{
		Role:     message.Role,
		Content:  message.Content,
		Metadata: metadata,
	}, nil
}

func (g *GPT4O) ConvertProviderResponseToMessage(response interface{}) (*models.Message, error) {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response is not a map")
	}

	content, ok := responseMap["content"].(string)
	if !ok {
		return nil, fmt.Errorf("response content is not a string")
	}

	metadata, ok := responseMap["metadata"].(map[string]string)
	if !ok {
		return nil, fmt.Errorf("response metadata is not a map")
	}

	role, ok := responseMap["role"].(string)
	if !ok {
		return nil, fmt.Errorf("response role is not a string")
	}

	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return &models.Message{
		Role:     role,
		Content:  content,
		Metadata: metadataJson,
	}, nil
}
