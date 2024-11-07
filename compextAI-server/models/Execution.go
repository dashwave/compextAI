package models

import (
	"encoding/json"
	"fmt"

	"github.com/burnerlee/compextAI/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ThreadExecutionStatus_IN_PROGRESS = "in_progress"
	ThreadExecutionStatus_COMPLETED   = "completed"
	ThreadExecutionStatus_FAILED      = "failed"
)

type ThreadExecution struct {
	Base
	ThreadID                  string                `json:"thread_id"`
	Thread                    Thread                `json:"thread" gorm:"foreignKey:ThreadID;references:Identifier"`
	ThreadExecutionParamID    string                `json:"thread_execution_param_id"`
	ThreadExecutionParams     ThreadExecutionParams `json:"thread_execution_params" gorm:"foreignKey:ThreadExecutionParamID;references:Identifier"`
	Status                    string                `json:"status"`
	Output                    json.RawMessage       `json:"output" gorm:"type:jsonb"`
	Content                   string                `json:"content"`
	Role                      string                `json:"role"`
	ExecutionResponseMetadata json.RawMessage       `json:"execution_response_metadata" gorm:"type:jsonb"`
	Metadata                  json.RawMessage       `json:"metadata" gorm:"type:jsonb"`
}

// ThreadExecutionParams are the parameters for executing a thread
type ThreadExecutionParams struct {
	Base
	UserID              uint            `json:"user_id"`
	Name                string          `json:"name"`
	Environment         string          `json:"environment"`
	Model               string          `json:"model"`
	Temperature         float64         `json:"temperature"`
	Timeout             int             `json:"timeout"`
	MaxTokens           int             `json:"max_tokens"`
	MaxCompletionTokens int             `json:"max_completion_tokens"`
	TopP                float64         `json:"top_p"`
	MaxOutputTokens     int             `json:"max_output_tokens"`
	ResponseFormat      json.RawMessage `json:"response_format" gorm:"type:jsonb"`
	SystemPrompt        string          `json:"system_prompt"`
}

func CreateThreadExecution(db *gorm.DB, threadExecution *ThreadExecution) (*ThreadExecution, error) {
	// create a new thread_execution_id
	threadExecutionIDUniqueIdentifier := uuid.New().String()
	threadExecutionID := fmt.Sprintf("%s%s", constants.THREAD_EXECUTION_ID_PREFIX, threadExecutionIDUniqueIdentifier)
	threadExecution.Identifier = threadExecutionID
	if err := db.Create(threadExecution).Error; err != nil {
		return nil, err
	}
	return threadExecution, nil
}

func UpdateThreadExecution(db *gorm.DB, threadExecution *ThreadExecution) error {
	updateData := make(map[string]interface{})
	if threadExecution.Status != "" {
		updateData["status"] = threadExecution.Status
	}
	if threadExecution.Output != nil {
		updateData["output"] = threadExecution.Output
	}
	if threadExecution.Content != "" {
		updateData["content"] = threadExecution.Content
	}
	if threadExecution.Role != "" {
		updateData["role"] = threadExecution.Role
	}
	if threadExecution.ExecutionResponseMetadata != nil {
		updateData["execution_response_metadata"] = threadExecution.ExecutionResponseMetadata
	}
	if threadExecution.Metadata != nil {
		updateData["metadata"] = threadExecution.Metadata
	}
	return db.Model(&ThreadExecution{}).Where("identifier = ?", threadExecution.Identifier).Updates(updateData).Error
}

func GetThreadExecutionByID(db *gorm.DB, executionID string) (*ThreadExecution, error) {
	var threadExecution ThreadExecution
	if err := db.Where("identifier = ?", executionID).Preload("Thread").First(&threadExecution).Error; err != nil {
		return nil, err
	}
	return &threadExecution, nil
}

func GetThreadExecutionParamsByID(db *gorm.DB, threadExecutionParamsID string) (*ThreadExecutionParams, error) {
	var threadExecutionParams ThreadExecutionParams
	if err := db.Where("identifier = ?", threadExecutionParamsID).First(&threadExecutionParams).Error; err != nil {
		return nil, err
	}
	return &threadExecutionParams, nil
}

func GetAllThreadExecutionParams(db *gorm.DB, userID uint) ([]ThreadExecutionParams, error) {
	var threadExecutionParams []ThreadExecutionParams
	if err := db.Where("user_id = ?", userID).Find(&threadExecutionParams).Error; err != nil {
		return nil, err
	}
	return threadExecutionParams, nil
}

func CreateThreadExecutionParams(db *gorm.DB, threadExecutionParams *ThreadExecutionParams) (*ThreadExecutionParams, error) {
	threadExecutionParamsIDUniqueIdentifier := uuid.New().String()
	threadExecutionParamsID := fmt.Sprintf("%s%s", constants.THREAD_EXECUTION_PARAMS_ID_PREFIX, threadExecutionParamsIDUniqueIdentifier)
	threadExecutionParams.Identifier = threadExecutionParamsID
	if err := db.Create(threadExecutionParams).Error; err != nil {
		return nil, err
	}
	return threadExecutionParams, nil
}

func UpdateThreadExecutionParams(db *gorm.DB, threadExecutionParams *ThreadExecutionParams) error {
	updateData := make(map[string]interface{})
	if threadExecutionParams.Model != "" {
		updateData["model"] = threadExecutionParams.Model
	}
	if threadExecutionParams.Temperature != 0 {
		updateData["temperature"] = threadExecutionParams.Temperature
	}
	if threadExecutionParams.Timeout != 0 {
		updateData["timeout"] = threadExecutionParams.Timeout
	}
	if threadExecutionParams.MaxTokens != 0 {
		updateData["max_tokens"] = threadExecutionParams.MaxTokens
	}
	if threadExecutionParams.MaxCompletionTokens != 0 {
		updateData["max_completion_tokens"] = threadExecutionParams.MaxCompletionTokens
	}
	if threadExecutionParams.MaxOutputTokens != 0 {
		updateData["max_output_tokens"] = threadExecutionParams.MaxOutputTokens
	}
	if threadExecutionParams.SystemPrompt != "" {
		updateData["system_prompt"] = threadExecutionParams.SystemPrompt
	}
	if threadExecutionParams.ResponseFormat != nil {
		updateData["response_format"] = threadExecutionParams.ResponseFormat
	}
	if threadExecutionParams.TopP != 0 {
		updateData["top_p"] = threadExecutionParams.TopP
	}
	return db.Model(&ThreadExecutionParams{}).Where("identifier = ?", threadExecutionParams.Identifier).Updates(updateData).Error
}

func DeleteThreadExecutionParams(db *gorm.DB, threadExecutionParamsID string) error {
	return db.Delete(&ThreadExecutionParams{}, "identifier = ?", threadExecutionParamsID).Error
}

func GetThreadExecutionParamsByUserIDAndNameAndEnvironment(db *gorm.DB, userID uint, name string, environment string) (*ThreadExecutionParams, error) {
	var threadExecutionParams ThreadExecutionParams
	if err := db.Where("user_id = ? AND name = ? AND environment = ?", userID, name, environment).First(&threadExecutionParams).Error; err != nil {
		return nil, err
	}
	return &threadExecutionParams, nil
}
