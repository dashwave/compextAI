package chat

import (
	"fmt"

	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

var (
	chatCompletionsProviderRegistry *ChatCompletionsProviderRegistry
)

type ChatCompletionsProvider interface {
	ValidateMessage(message *models.Message) error
	ConvertMessageToProviderFormat(message *models.Message) (interface{}, error)
	ConvertExecutionResponseToMessage(response interface{}) (*models.Message, error)
	GetProviderOwner() string
	GetProviderModel() string
	GetProviderIdentifier() string
	ExecuteThread(db *gorm.DB, user *models.User, messages []*models.Message, threadExecutionParams *models.ThreadExecutionParams, threadExecutionIdentifier string) (int, interface{}, error)
}

type ChatCompletionsProvider_Enum string

type ChatCompletionsProviderRegistry struct {
	providers map[ChatCompletionsProvider_Enum]ChatCompletionsProvider
}

func NewChatCompletionsProviderRegistry() *ChatCompletionsProviderRegistry {
	return &ChatCompletionsProviderRegistry{
		providers: make(map[ChatCompletionsProvider_Enum]ChatCompletionsProvider),
	}
}

func (r *ChatCompletionsProviderRegistry) register(provider ChatCompletionsProvider) {
	providerIdentifier := ChatCompletionsProvider_Enum(provider.GetProviderIdentifier())
	r.providers[providerIdentifier] = provider
}

func GetChatCompletionsProvider(providerIdentifier string) (ChatCompletionsProvider, error) {
	provider, ok := chatCompletionsProviderRegistry.providers[ChatCompletionsProvider_Enum(providerIdentifier)]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", providerIdentifier)
	}
	return provider, nil
}

func (r *ChatCompletionsProviderRegistry) getAvailableProviders() []ChatCompletionsProvider {
	providers := make([]ChatCompletionsProvider, 0)
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers
}
