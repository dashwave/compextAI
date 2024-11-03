package models

type Thread struct {
	Base
	Title string `json:"title" gorm:"not null"`
	Metadata map[string]string `json:"metadata"`
}
