package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Message struct {
	Base
	Content  string          `json:"content" gorm:"not null"`
	Role     string          `json:"role" gorm:"not null"`
	ThreadID string          `json:"thread_id" gorm:"not null"`
	Thread   Thread          `json:"thread" gorm:"foreignKey:ThreadID;references:Identifier"`
	Metadata json.RawMessage `json:"metadata"`

	// Implement support for tool calls and function calls later on
	// ToolCalls []ToolCall        `json:"tool_calls"`
	// FunctionCall FunctionCall      `json:"function_call"`
}

func GetAllMessages(db *gorm.DB, threadID string) ([]Message, error) {
	var messages []Message
	if err := db.Where("thread_id = ?", threadID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func CreateMessage(db *gorm.DB, message *Message) error {
	return db.Create(message).Error
}

func GetMessage(db *gorm.DB, messageID string) (*Message, error) {
	var message Message
	if err := db.First(&message, "identifier = ?", messageID).Preload("Thread").Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func UpdateMessage(db *gorm.DB, message *Message) (*Message, error) {
	updateData := make(map[string]interface{})
	if message.Role != "" {
		updateData["role"] = message.Role
	}
	if message.Metadata != nil {
		updateData["metadata"] = message.Metadata
	}
	if message.Content != "" {
		updateData["content"] = message.Content
	}
	if err := db.Model(&Message{}).Where("identifier = ?", message.Identifier).Updates(updateData).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func DeleteMessage(db *gorm.DB, messageID string) error {
	return db.Where("identifier = ?", messageID).Delete(&Message{}).Error
}
