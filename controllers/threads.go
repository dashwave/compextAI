package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/burnerlee/compextAI/constants"
	"github.com/burnerlee/compextAI/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateThread(db *gorm.DB, request *CreateThreadRequest) (*models.Thread, error) {
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// create a new thread_id
	threadIDUniqueIdentifier := uuid.New().String()
	threadID := fmt.Sprintf("%s%s", constants.THREAD_ID_PREFIX, threadIDUniqueIdentifier)

	metadataJsonBlob, err := json.Marshal(request.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	thread := models.Thread{
		UserID: request.UserID,
		Base: models.Base{
			Identifier: threadID,
		},
		Title:    request.Title,
		Metadata: metadataJsonBlob,
	}

	if err := models.CreateThread(tx, &thread); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create thread: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &thread, nil
}
