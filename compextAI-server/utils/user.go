package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/burnerlee/compextAI/models"
	"gorm.io/gorm"
)

func GetUserIDFromRequest(r *http.Request) (int, error) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return 0, errors.New("user ID not found")
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return 0, err
	}

	return userIDInt, nil
}

func CheckThreadAccess(db *gorm.DB, threadID string, userID uint) (bool, error) {
	thread, err := models.GetThread(db, threadID)
	if err != nil {
		return false, err
	}

	return thread.UserID == userID, nil
}

func CheckMessageAccess(db *gorm.DB, messageID string, userID uint) (bool, error) {
	message, err := models.GetMessage(db, messageID)
	if err != nil {
		return false, err
	}

	return message.Thread.UserID == userID, nil
}

func CheckThreadExecutionAccess(db *gorm.DB, executionID string, userID uint) (bool, error) {
	threadExecution, err := models.GetThreadExecutionByID(db, executionID)
	if err != nil {
		return false, err
	}

	return threadExecution.UserID == userID, nil
}
