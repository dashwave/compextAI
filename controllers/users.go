package controllers

import (
	"errors"

	"github.com/burnerlee/compextAI/models"
	"github.com/burnerlee/compextAI/utils"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, request *CreateUserRequest) (*models.User, error) {

	user := &models.User{
		Username: request.Username,
		Password: request.Password,
	}

	apiToken, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	user.APIToken = apiToken

	tx := db.Begin()
	if err := models.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

func Login(db *gorm.DB, request *LoginRequest) (*models.User, error) {
	user, err := models.GetUserByUsername(db, request.Username)
	if err != nil {
		return nil, err
	}

	if user.Password != request.Password {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
