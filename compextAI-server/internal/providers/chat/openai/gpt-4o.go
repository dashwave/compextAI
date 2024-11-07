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
	Role     string                 `json:"role"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (g *GPT4O) ConvertMessageToProviderFormat(message *models.Message) (interface{}, error) {
	var metadata map[string]interface{}
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

func (g *GPT4O) ConvertExecutionResponseToMessage(response interface{}) (*models.Message, error) {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response is not a map")
	}

	messageChoices := responseMap["choices"].([]interface{})
	if len(messageChoices) == 0 {
		return nil, fmt.Errorf("no message choices found")
	}
	messageChoice := messageChoices[0].(map[string]interface{})
	message := messageChoice["message"].(map[string]interface{})

	role, ok := message["role"].(string)
	if !ok {
		return nil, fmt.Errorf("message role is not a string")
	}
	content, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("message content is not a string")
	}

	openAIChatCompletionID := responseMap["id"].(string)
	usage := responseMap["usage"].(map[string]interface{})

	metadata := map[string]interface{}{
		"openai_chat_completion_id": openAIChatCompletionID,
		"usage":                     usage,
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

func (g *GPT4O) ExecuteThread(db *gorm.DB, user *models.User, messages []*models.Message, threadExecutionParams *models.ThreadExecutionParams, threadExecutionIdentifier string) (int, interface{}, error) {
	systemPrompt := ""

	modelMessages := make([]gpt4oOpenAIMessage, 0)
	for _, message := range messages {
		modelMessage, err := g.ConvertMessageToProviderFormat(message)
		if err != nil {
			logger.GetLogger().Errorf("Error converting message to provider format: %v", err)
			return -1, nil, err
		}
		if message.Role == "system" {
			systemPrompt = message.Content
			continue
		}
		modelMessages = append(modelMessages, modelMessage.(gpt4oOpenAIMessage))
	}

	// override the system prompt if it is provided for execution
	if threadExecutionParams.Template.SystemPrompt != "" {
		systemPrompt = threadExecutionParams.Template.SystemPrompt
	}

	// add the system prompt to the beginning of the messages thread if it is provided
	if systemPrompt != "" {
		modelMessages = append([]gpt4oOpenAIMessage{{
			Role:    "system",
			Content: systemPrompt,
		}}, modelMessages...)
	}

	if threadExecutionParams.Template.Temperature == 0 {
		threadExecutionParams.Template.Temperature = DEFAULT_TEMPERATURE
	}
	if threadExecutionParams.Template.MaxCompletionTokens == 0 {
		threadExecutionParams.Template.MaxCompletionTokens = DEFAULT_MAX_COMPLETION_TOKENS
	}
	if threadExecutionParams.Template.Timeout == 0 {
		threadExecutionParams.Template.Timeout = DEFAULT_TIMEOUT
	}

	executionData := gpt4oExecutionData{
		APIKey:              user.OpenAIKey,
		Model:               g.model,
		Messages:            modelMessages,
		Temperature:         threadExecutionParams.Template.Temperature,
		MaxCompletionTokens: threadExecutionParams.Template.MaxCompletionTokens,
		Timeout:             threadExecutionParams.Template.Timeout,
		ResponseFormat:      threadExecutionParams.Template.ResponseFormat,
	}

	if err := executionData.Validate(); err != nil {
		logger.GetLogger().Errorf("Error validating execution data: %v", err)
		return -1, nil, err
	}

	if err := base.UpdateThreadExecutionMetadata(db, threadExecutionIdentifier, executionData); err != nil {
		logger.GetLogger().Errorf("Error updating thread execution metadata: %v", err)
		return -1, nil, err
	}

	executionParams := &base.ExecuteParams{
		Timeout: time.Duration(executionData.Timeout) * time.Second,
	}

	return base.Execute(g.executorRoute, executionParams, executionData)
}
