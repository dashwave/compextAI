package models

import "gorm.io/gorm"

type User struct {
	Base
	Username     string `json:"username" gorm:"unique"`
	Password     string `json:"password" gorm:"not null"`
	APIToken     string `json:"api_token" gorm:"unique"`
	OpenAIKey    string `json:"openai_key"`
	AnthropicKey string `json:"anthropic_key"`
}

func GetUserIDByAPIToken(db *gorm.DB, token string) (uint, error) {
	var user User
	if err := db.Where("api_token = ?", token).First(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
