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
	ThreadID              string          `json:"thread_id"`
	Thread                Thread          `json:"thread" gorm:"foreignKey:ThreadID;references:Identifier"`
	ThreadExecutionParams json.RawMessage `json:"thread_execution_params"`
	Status                string          `json:"status"`
	Output                string          `json:"output"`
	ResponseContent       string          `json:"response_content"`
	ResponseRole          string          `json:"response_role"`
}

// ThreadExecutionParams are the parameters for executing a thread
type ThreadExecutionParams struct {
	Model               string      `json:"model"`
	Temperature         float64     `json:"temperature"`
	Timeout             int         `json:"timeout"`
	MaxCompletionTokens int         `json:"max_completion_tokens"`
	TopP                float64     `json:"top_p"`
	MaxOutputTokens     int         `json:"max_output_tokens"`
	ResponseFormat      interface{} `json:"response_format"`
}

func CreateThreadExecution(db *gorm.DB, threadExecution *ThreadExecution) error {
	// create a new thread_execution_id
	threadExecutionIDUniqueIdentifier := uuid.New().String()
	threadExecutionID := fmt.Sprintf("%s%s", constants.THREAD_EXECUTION_ID_PREFIX, threadExecutionIDUniqueIdentifier)
	threadExecution.Identifier = threadExecutionID
	return db.Create(threadExecution).Error
}

func UpdateThreadExecution(db *gorm.DB, threadExecution *ThreadExecution) error {
	updateData := make(map[string]interface{})
	if threadExecution.Status != "" {
		updateData["status"] = threadExecution.Status
	}
	if threadExecution.Output != "" {
		updateData["output"] = threadExecution.Output
	}
	if threadExecution.ResponseContent != "" {
		updateData["response_content"] = threadExecution.ResponseContent
	}
	if threadExecution.ResponseRole != "" {
		updateData["response_role"] = threadExecution.ResponseRole
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
