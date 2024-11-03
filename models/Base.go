package models

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string         `json:"id" gorm:"primary_key;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
