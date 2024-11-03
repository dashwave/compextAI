package models

import "gorm.io/gorm"

type Message struct {
	Base
	Content string `json:"content" gorm:"not null"`
	Role    string `json:"role" gorm:"not null"`
	ThreadID string `json:"thread_id" gorm:"not null;index"`
	Thread Thread `json:"thread" gorm:"foreignKey:ThreadID"`
	Metadata map[string]string `json:"metadata"`
}
