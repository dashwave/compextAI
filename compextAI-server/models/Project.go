package models

import (
	"fmt"

	"github.com/burnerlee/compextAI/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	Base
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateProject(db *gorm.DB, project *Project) error {
	projectID := uuid.New().String()
	projectIdentifier := fmt.Sprintf("%s%s", constants.PROJECT_ID_PREFIX, projectID)
	project.Identifier = projectIdentifier
	return db.Create(project).Error
}

func GetProject(db *gorm.DB, projectID string) (*Project, error) {
	var project Project
	if err := db.First(&project, "identifier = ?", projectID).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func GetProjectByName(db *gorm.DB, name string, userID uint) (*Project, error) {
	var project Project
	if err := db.First(&project, "name = ? AND user_id = ?", name, userID).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func GetAllProjects(db *gorm.DB, userID uint) ([]Project, error) {
	var projects []Project
	if err := db.Where("user_id = ?", userID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func DeleteProject(db *gorm.DB, projectID string) error {
	return db.Delete(&Project{}, "identifier = ?", projectID).Error
}

func UpdateProject(db *gorm.DB, project *Project) error {
	var updateData = make(map[string]interface{})
	if project.Name != "" {
		updateData["name"] = project.Name
	}
	if project.Description != "" {
		updateData["description"] = project.Description
	}
	return db.Model(&Project{}).Where("identifier = ?", project.Identifier).Updates(updateData).Error
}
