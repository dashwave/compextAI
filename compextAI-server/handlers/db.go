package handlers

import (
	"errors"
	"fmt"
	"os"

	"github.com/burnerlee/compextAI/constants"
	"github.com/burnerlee/compextAI/internal/logger"
	"github.com/burnerlee/compextAI/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadEnv() {
	fmt.Println("Loading environment variables")
	err := godotenv.Load(".env")
	if err != nil {
		logger.GetLogger().Warnf("Error loading .env file: %v", err)
	} else {
		logger.GetLogger().Info("Environment variables loaded")
	}
}

func InitDB() (*gorm.DB, error) {
	loadEnv()

	DbHost := os.Getenv("POSTGRES_HOST")
	DbPort := os.Getenv("POSTGRES_PORT")
	DbUser := os.Getenv("POSTGRES_USER")
	DbName := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	DbPassword := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", DbHost, DbPort, DbUser, DbName, sslMode, DbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Message{}, &models.Thread{}, &models.User{}, &models.ThreadExecution{}, &models.ThreadExecutionParams{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	adminUser, err := models.GetUserByUsername(db, "admin")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			adminUser = &models.User{
				Base: models.Base{
					Identifier: "admin",
				},
				Username: "admin",
			}
			if err := db.Create(adminUser).Error; err != nil {
				return fmt.Errorf("failed to create admin user: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get admin user: %w", err)
		}
	}

	_, err = models.GetThread(db, constants.THREAD_IDENTIFIER_FOR_NULL_THREAD)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			nullThread := &models.Thread{
				Base: models.Base{
					Identifier: constants.THREAD_IDENTIFIER_FOR_NULL_THREAD,
				},
				UserID: adminUser.ID,
				User:   *adminUser,
			}
			if err := db.Create(nullThread).Error; err != nil {
				return fmt.Errorf("failed to create null thread: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get null thread: %w", err)
		}
	}

	return nil
}
