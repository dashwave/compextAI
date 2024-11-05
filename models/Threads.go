package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// Thread is a collection of messages
type Thread struct {
	Base
	User     User            `json:"user" gorm:"foreignKey:UserID"`
	UserID   uint            `json:"user_id" gorm:"not null"`
	Title    string          `json:"title" gorm:"not null"`
	Metadata json.RawMessage `json:"metadata"`
}

func GetAllThreads(db *gorm.DB, userID uint) ([]Thread, error) {
	var threads []Thread
	if err := db.Where("user_id = ?", userID).Find(&threads).Error; err != nil {
		return nil, err
	}
	return threads, nil
}

func CreateThread(db *gorm.DB, thread *Thread) error {
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
