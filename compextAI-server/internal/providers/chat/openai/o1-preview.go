package openai

import (
	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

const (
	O1_PREVIEW_MODEL          = "o1-preview"
	O1_PREVIEW_OWNER          = "openai"
	O1_PREVIEW_IDENTIFIER     = "o1-preview"
	O1_PREVIEW_EXECUTOR_ROUTE = "/chatcompletion/openai"

	O1_PREVIEW_DEFAULT_TEMPERATURE           = 1
	O1_PREVIEW_DEFAULT_MAX_COMPLETION_TOKENS = 32768
	O1_PREVIEW_DEFAULT_TIMEOUT               = 600
)

type O1Preview struct {
	owner         string
	model         string
	allowedRoles  []string
	executorRoute string
}

func NewO1Preview() *O1Preview {
	return &O1Preview{
		owner:         O1_PREVIEW_OWNER,
		model:         O1_PREVIEW_MODEL,
		allowedRoles:  openaiAllowedRoles,
		executorRoute: O1_PREVIEW_EXECUTOR_ROUTE,
	}
}

func (g *O1Preview) GetProviderOwner() string {
	return g.owner
}

func (g *O1Preview) GetProviderModel() string {
	return g.model
}

func (g *O1Preview) GetProviderIdentifier() string {
	return O1_PREVIEW_IDENTIFIER
}

func (g *O1Preview) ValidateMessage(message *models.Message) error {
	return validateMessage(message)
}

func (g *O1Preview) ConvertMessageToProviderFormat(message *models.Message) (interface{}, error) {
	return convertMessageToProviderFormat(message)
}

func (g *O1Preview) ConvertExecutionResponseToMessage(response interface{}) (*models.Message, error) {
	return convertExecutionResponseToMessage(response)
}

func (g *O1Preview) ExecuteThread(db *gorm.DB, user *models.User, messages []*models.Message, threadExecutionParamsTemplate *models.ThreadExecutionParamsTemplate, threadExecutionIdentifier string) (int, interface{}, error) {
	return executeThread(db, user, messages, threadExecutionParamsTemplate, threadExecutionIdentifier, &executeParamConfigs{
		Model:                      g.model,
		ExecutorRoute:              g.executorRoute,
		DefaultTemperature:         O1_PREVIEW_DEFAULT_TEMPERATURE,
		DefaultMaxCompletionTokens: O1_PREVIEW_DEFAULT_MAX_COMPLETION_TOKENS,
		DefaultTimeout:             O1_PREVIEW_DEFAULT_TIMEOUT,
	})
}
