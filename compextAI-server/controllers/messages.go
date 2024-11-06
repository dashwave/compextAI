package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

func CreateMessage(db *gorm.DB, req *CreateMessageRequest) (*models.Message, error) {
	// validate if the thread exists
	if _, err := models.GetThread(db, req.ThreadID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("thread not found")
		}
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	metadataJsonBlob, err := json.Marshal(req.Metadata)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	message := &models.Message{
		ThreadID: req.ThreadID,
		Content:  req.Content,
		Role:     req.Role,
		Metadata: metadataJsonBlob,
	}

	if err := models.CreateMessage(tx, message); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	tx.Commit()
	return message, nil
}
