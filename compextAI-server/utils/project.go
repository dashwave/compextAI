package utils

import (
	"errors"

	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

func GetProjectIDFromName(db *gorm.DB, name string, userID uint) (string, error) {
	project, err := models.GetProjectByName(db, name, userID)
	if err != nil {
		return "", err
	}

	hasAccess, err := CheckProjectAccess(db, project.Identifier, userID)
	if err != nil {
		return "", err
	}
	if !hasAccess {
		return "", errors.New("you do not have access to this project")
	}

	return project.Identifier, nil
}
