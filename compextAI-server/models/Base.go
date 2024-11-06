package models

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID         uint           `json:"id" gorm:"primary_key"`
	Identifier string         `json:"identifier" gorm:"unique;index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}
