package models

import (
	"encoding/json"
	"fmt"

	"github.com/burnerlee/compextAI/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Thread is a collection of messages
type Thread struct {
	Base
	User      User            `json:"user" gorm:"foreignKey:UserID"`
	UserID    uint            `json:"user_id" gorm:"not null"`
	ProjectID string          `json:"project_id"`
	Title     string          `json:"title" gorm:"not null"`
	Metadata  json.RawMessage `json:"metadata" gorm:"type:jsonb;default:'{}'"`
}

func GetAllThreads(db *gorm.DB, userID uint, projectID string, searchQuery string, searchFiltersMap map[string]string, page, limit int) ([]Thread, int64, error) {
	offset := (page - 1) * limit
	var total int64

	query := db.Model(&Thread{}).Where("user_id = ? AND project_id = ?", userID, projectID)

	if searchQuery != "" {
		query = query.Where("title LIKE ? OR identifier LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	if len(searchFiltersMap) > 0 {
		for key, value := range searchFiltersMap {
			// apply the search filters on metadata keys
			query = query.Where("metadata ->> ? = ?", key, value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC").Offset(offset).Limit(limit)

	var threads []Thread
	if err := query.Find(&threads).Error; err != nil {
		return nil, 0, err
	}

	return threads, total, nil
}

func CreateThread(db *gorm.DB, thread *Thread) error {
	// create a new thread_id
	threadIDUniqueIdentifier := uuid.New().String()
	threadID := fmt.Sprintf("%s%s", constants.THREAD_ID_PREFIX, threadIDUniqueIdentifier)
	thread.Identifier = threadID
	return db.Create(thread).Error
}

func GetThread(db *gorm.DB, threadID string) (*Thread, error) {
	var thread Thread
	if err := db.Where("identifier = ?", threadID).First(&thread).Error; err != nil {
		return nil, err
	}
	return &thread, nil
}

func UpdateThread(db *gorm.DB, thread *Thread) (*Thread, error) {
	// update the thread in the db
	updateData := make(map[string]interface{})
	if thread.Title != "" {
		updateData["title"] = thread.Title
	}
	if thread.Metadata != nil {
		updateData["metadata"] = thread.Metadata
	}
	if err := db.Model(&Thread{}).Where("identifier = ?", thread.Identifier).Updates(updateData).Error; err != nil {
		return nil, err
	}
	return thread, nil
}

func DeleteThread(db *gorm.DB, threadID string) error {
	return db.Where("identifier = ?", threadID).Delete(&Thread{}).Error
}

func (t *Thread) GetAllMessages(db *gorm.DB) ([]Message, error) {
	var messages []Message
	if err := db.Where("thread_id = ?", t.Identifier).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
