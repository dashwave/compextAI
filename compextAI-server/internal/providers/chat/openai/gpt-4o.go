package openai

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/burnerlee/compextAI/internal/logger"
	"github.com/burnerlee/compextAI/internal/providers/chat/base"
	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

const (
	GPT4O_MODEL          = "gpt-4o"
	GPT4O_OWNER          = "openai"
	GPT4O_IDENTIFIER     = "gpt-4o"
	GPT4O_EXECUTOR_ROUTE = "/chatcompletion/openai"

	DEFAULT_TEMPERATURE           = 0.5
	DEFAULT_MAX_COMPLETION_TOKENS = 10000
	DEFAULT_TIMEOUT               = 600
)

type GPT4O struct {
	owner         string
	model         string
	allowedRoles  []string
	executorRoute string
}

func NewGPT4O() *GPT4O {
	return &GPT4O{
		owner:         GPT4O_OWNER,
		model:         GPT4O_MODEL,
		allowedRoles:  []string{"user", "assistant", "system"},
		executorRoute: GPT4O_EXECUTOR_ROUTE,
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
	if message.Metadata != nil {
		if err := json.Unmarshal(message.Metadata, &metadata); err != nil {
			return nil, err
		}
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

type gpt4oExecutionData struct {
	APIKey              string               `json:"api_key"`
	Model               string               `json:"model"`
	Messages            []gpt4oOpenAIMessage `json:"messages"`
	Temperature         float64              `json:"temperature"`
	MaxCompletionTokens int                  `json:"max_completion_tokens"`
	Timeout             int                  `json:"timeout"`
	ResponseFormat      interface{}          `json:"response_format"`
}

func (d *gpt4oExecutionData) Validate() error {
	return nil
}

func (g *GPT4O) ExecuteThread(db *gorm.DB, user *models.User, thread *models.Thread, threadExecutionParams *models.ThreadExecutionParams) (int, interface{}, error) {
	threadMessages, err := thread.GetAllMessages(db)
	if err != nil {
		logger.GetLogger().Errorf("Error getting thread messages: %v", err)
		return -1, nil, err
	}

	modelMessages := make([]gpt4oOpenAIMessage, 0)
	for _, message := range threadMessages {
		modelMessage, err := g.ConvertMessageToProviderFormat(&message)
		if err != nil {
			logger.GetLogger().Errorf("Error converting message to provider format: %v", err)
			return -1, nil, err
		}
		modelMessages = append(modelMessages, modelMessage.(gpt4oOpenAIMessage))
	}

	if threadExecutionParams.Temperature == 0 {
		threadExecutionParams.Temperature = DEFAULT_TEMPERATURE
	}
	if threadExecutionParams.MaxCompletionTokens == 0 {
		threadExecutionParams.MaxCompletionTokens = DEFAULT_MAX_COMPLETION_TOKENS
	}
	if threadExecutionParams.Timeout == 0 {
		threadExecutionParams.Timeout = DEFAULT_TIMEOUT
	}

	executionData := gpt4oExecutionData{
		APIKey:              user.OpenAIKey,
		Model:               g.model,
		Messages:            modelMessages,
		Temperature:         threadExecutionParams.Temperature,
		MaxCompletionTokens: threadExecutionParams.MaxCompletionTokens,
		Timeout:             threadExecutionParams.Timeout,
		ResponseFormat:      threadExecutionParams.ResponseFormat,
	}

	if err := executionData.Validate(); err != nil {
		logger.GetLogger().Errorf("Error validating execution data: %v", err)
		return -1, nil, err
	}

	executionParams := &base.ExecuteParams{
		Timeout: time.Duration(executionData.Timeout) * time.Second,
	}

	return base.Execute(g.executorRoute, executionParams, executionData)
}
