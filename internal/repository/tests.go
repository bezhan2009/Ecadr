package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
)

func GetTasksByTypeAndID(targetType string, targetID uint) (tests []models.Test, err error) {
	if err = db.GetDBConn().Model(&models.Test{}).Where("target_type = ? AND target_id = ?", targetType, targetID).Find(&tests).Error; err != nil {
		logger.Error.Printf("[repository.GetTasksByTypeAndID] Error getting tasks by type and id %v", err)

		return nil, err
	}

	return tests, nil
}

func GetTaskByID(taskID uint) (test *models.Test, err error) {
	if err = db.GetDBConn().Model(&models.Test{}).Where("id = ?", taskID).Find(&test).Error; err != nil {
		logger.Error.Printf("[repository.GetTaskByID] Error getting task by id %v", err)
		return nil, err
	}

	return test, nil
}
