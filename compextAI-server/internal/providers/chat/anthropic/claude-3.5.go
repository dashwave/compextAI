package anthropic

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
	ANTHROPIC_MODEL          = "claude-3-5-sonnet-20241022"
	ANTHROPIC_OWNER          = "anthropic"
	ANTHROPIC_IDENTIFIER     = "claude-3-5-sonnet"
	ANTHROPIC_EXECUTOR_ROUTE = "/chatcompletion/anthropic"

	DEFAULT_TEMPERATURE = 0.5
	DEFAULT_MAX_TOKENS  = 8192
	DEFAULT_TIMEOUT     = 600
)

type Claude35 struct {
	owner         string
	model         string
	allowedRoles  []string
	executorRoute string
}

func NewClaude35() *Claude35 {
	return &Claude35{
		owner:         ANTHROPIC_OWNER,
		model:         ANTHROPIC_MODEL,
		allowedRoles:  []string{"user", "assistant", "system"},
		executorRoute: ANTHROPIC_EXECUTOR_ROUTE,
	}
}

func (g *Claude35) GetProviderOwner() string {
	return g.owner
}

func (g *Claude35) GetProviderModel() string {
	return g.model
}

func (g *Claude35) GetProviderIdentifier() string {
	return ANTHROPIC_IDENTIFIER
}

func (g *Claude35) ValidateMessage(message *models.Message) error {
	if message.Content == "" {
		return fmt.Errorf("message content is empty")
	}

	if !slices.Contains(g.allowedRoles, message.Role) {
		return fmt.Errorf("message role is invalid, only %v are allowed", g.allowedRoles)
	}
	return nil
}

type claude35Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (g *Claude35) ConvertMessageToProviderFormat(message *models.Message) (interface{}, error) {

	return claude35Message{
		Role:    message.Role,
		Content: message.Content,
	}, nil
}

func (g *Claude35) ConvertExecutionResponseToMessage(response interface{}) (*models.Message, error) {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response is not a map")
	}

	contentChoices := responseMap["content"].([]interface{})
	if len(contentChoices) == 0 {
		return nil, fmt.Errorf("no content found")
	}
	contentChoice := contentChoices[0].(map[string]interface{})

	content, ok := contentChoice["text"].(string)
	if !ok {
		return nil, fmt.Errorf("content is not a string")
	}

	role, ok := responseMap["role"].(string)
	if !ok {
		return nil, fmt.Errorf("role is not a string")
	}

	anthropicChatCompletionID := responseMap["id"].(string)
	usage := responseMap["usage"].(map[string]interface{})

	metadata := map[string]interface{}{
		"anthropic_chat_completion_id": anthropicChatCompletionID,
		"usage":                        usage,
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

type claude35ExecutionData struct {
	APIKey       string            `json:"api_key"`
	Model        string            `json:"model"`
	Messages     []claude35Message `json:"messages"`
	Temperature  float64           `json:"temperature"`
	Timeout      int               `json:"timeout"`
	MaxTokens    int               `json:"max_tokens"`
	SystemPrompt string            `json:"system_prompt"`
}

func (d *claude35ExecutionData) Validate() error {
	return nil
}

func (g *Claude35) ExecuteThread(db *gorm.DB, user *models.User, messages []*models.Message, threadExecutionParams *models.ThreadExecutionParams, threadExecutionIdentifier string) (int, interface{}, error) {
	systemPrompt := ""

	modelMessages := make([]claude35Message, 0)
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
		modelMessages = append(modelMessages, modelMessage.(claude35Message))
	}

	// override the system prompt if it is provided for execution
	if threadExecutionParams.SystemPrompt != "" {
		systemPrompt = threadExecutionParams.SystemPrompt
	}

	if threadExecutionParams.Temperature == 0 {
		threadExecutionParams.Temperature = DEFAULT_TEMPERATURE
	}
	if threadExecutionParams.MaxTokens == 0 {
		threadExecutionParams.MaxTokens = DEFAULT_MAX_TOKENS
	}
	if threadExecutionParams.Timeout == 0 {
		threadExecutionParams.Timeout = DEFAULT_TIMEOUT
	}

	executionData := claude35ExecutionData{
		APIKey:       user.AnthropicKey,
		Model:        g.model,
		Messages:     modelMessages,
		Temperature:  threadExecutionParams.Temperature,
		MaxTokens:    threadExecutionParams.MaxTokens,
		Timeout:      threadExecutionParams.Timeout,
		SystemPrompt: systemPrompt,
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
