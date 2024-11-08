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

var (
	openaiAllowedRoles = []string{"user", "assistant", "system"}
)

func validateMessage(message *models.Message) error {
	if message.Content == "" {
		return fmt.Errorf("message content is empty")
	}

	if !slices.Contains(openaiAllowedRoles, message.Role) {
		return fmt.Errorf("message role is invalid, only %v are allowed", openaiAllowedRoles)
	}
	return nil
}

type openaiMessage struct {
	Role     string                 `json:"role"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

func convertMessageToProviderFormat(message *models.Message) (interface{}, error) {
	var metadata map[string]interface{}
	if message.Metadata != nil {
		if err := json.Unmarshal(message.Metadata, &metadata); err != nil {
			return nil, err
		}
	}

	return openaiMessage{
		Role:     message.Role,
		Content:  message.Content,
		Metadata: metadata,
	}, nil
}

func convertExecutionResponseToMessage(response interface{}) (*models.Message, error) {
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

type openaiExecutionData struct {
	APIKey              string          `json:"api_key"`
	Model               string          `json:"model"`
	Messages            []openaiMessage `json:"messages"`
	Temperature         float64         `json:"temperature"`
	MaxCompletionTokens int             `json:"max_completion_tokens"`
	Timeout             int             `json:"timeout"`
	ResponseFormat      interface{}     `json:"response_format"`
}

func (d *openaiExecutionData) Validate() error {
	return nil
}

type executeParamConfigs struct {
	Model                      string
	ExecutorRoute              string
	DefaultTemperature         float64
	DefaultMaxCompletionTokens int
	DefaultTimeout             int
}

func executeThread(db *gorm.DB, user *models.User, messages []*models.Message, threadExecutionParamsTemplate *models.ThreadExecutionParamsTemplate, threadExecutionIdentifier string, configs *executeParamConfigs) (int, interface{}, error) {
	systemPrompt := ""

	modelMessages := make([]openaiMessage, 0)
	for _, message := range messages {
		modelMessage, err := convertMessageToProviderFormat(message)
		if err != nil {
			logger.GetLogger().Errorf("Error converting message to provider format: %v", err)
			return -1, nil, err
		}
		if message.Role == "system" {
			systemPrompt = message.Content
			continue
		}
		modelMessages = append(modelMessages, modelMessage.(openaiMessage))
	}

	// override the system prompt if it is provided for execution
	if threadExecutionParamsTemplate.SystemPrompt != "" {
		systemPrompt = threadExecutionParamsTemplate.SystemPrompt
	}

	// add the system prompt to the beginning of the messages thread if it is provided
	if systemPrompt != "" {
		modelMessages = append([]openaiMessage{{
			Role:    "system",
			Content: systemPrompt,
		}}, modelMessages...)
	}

	if threadExecutionParamsTemplate.Temperature == 0 {
		threadExecutionParamsTemplate.Temperature = configs.DefaultTemperature
	}
	if threadExecutionParamsTemplate.MaxCompletionTokens == 0 {
		threadExecutionParamsTemplate.MaxCompletionTokens = configs.DefaultMaxCompletionTokens
	}
	if threadExecutionParamsTemplate.Timeout == 0 {
		threadExecutionParamsTemplate.Timeout = configs.DefaultTimeout
	}

	executionData := openaiExecutionData{
		APIKey:              user.OpenAIKey,
		Model:               configs.Model,
		Messages:            modelMessages,
		Temperature:         threadExecutionParamsTemplate.Temperature,
		MaxCompletionTokens: threadExecutionParamsTemplate.MaxCompletionTokens,
		Timeout:             threadExecutionParamsTemplate.Timeout,
		ResponseFormat:      threadExecutionParamsTemplate.ResponseFormat,
	}

	if err := executionData.Validate(); err != nil {
		logger.GetLogger().Errorf("Error validating execution data: %v", err)
		return -1, nil, err
	}

	if err := base.UpdateThreadExecutionMetadata(db, threadExecutionIdentifier, executionData, messages); err != nil {
		logger.GetLogger().Errorf("Error updating thread execution metadata: %v", err)
		return -1, nil, err
	}

	executionParams := &base.ExecuteParams{
		Timeout: time.Duration(executionData.Timeout) * time.Second,
	}

	return base.Execute(configs.ExecutorRoute, executionParams, executionData)
}
