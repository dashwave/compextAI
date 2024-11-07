package base

import (
	"encoding/json"
	"fmt"

	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

func UpdateThreadExecutionMetadata(db *gorm.DB, threadExecutionIdentifier string, metadata interface{}) error {
	threadExecution, err := models.GetThreadExecutionByID(db, threadExecutionIdentifier)
	if err != nil {
		return fmt.Errorf("error getting thread execution: %v", err)
	}

	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("error marshalling metadata: %v", err)
	}

	threadExecution.Metadata = metadataJson
	if err := models.UpdateThreadExecution(db, threadExecution); err != nil {
		return fmt.Errorf("error updating thread execution: %v", err)
	}

	return nil
}
