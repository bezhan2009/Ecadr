package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
)

func GetTestsByTypeAndID(targetType, targetID uint) (tests []models.Test, err error) {
	if err = db.GetDBConn().
		Preload("Questions").
		Preload("Questions.Choices").
		Model(&models.Test{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).Find(&tests).Error; err != nil {
		logger.Error.Printf("[repository.GetTestsByTypeAndID] Error getting tasks by type and id %v", err)

		return nil, TranslateGormError(err)
	}

	return tests, nil
}

func GetTestByID(testID uint) (test models.Test, err error) {
	if err = db.GetDBConn().
		Preload("Questions").
		Preload("Questions.Choices").
		Model(&models.Test{}).
		Where("id = ?", testID).First(&test).Error; err != nil {
		logger.Error.Printf("[repository.GetTestByID] Error getting test by id %v", err)
		return test, TranslateGormError(err)
	}

	return test, nil
}

func CreateTest(test models.Test) (testID uint, err error) {
	if err = db.GetDBConn().Model(&models.Test{}).Create(&test).Error; err != nil {
		logger.Error.Printf("[repository.CreateTest] Error creating test %v", err)
		return 0, TranslateGormError(err)
	}

	return uint(test.ID), err
}

func UpdateTest(test models.Test) (testID uint, err error) {
	if err = db.GetDBConn().Model(&models.Test{}).Where("id = ?", test.ID).Updates(&test).Error; err != nil {
		logger.Error.Printf("[repository.UpdateTest] Error updating test %v", err)
		return 0, TranslateGormError(err)
	}

	return uint(test.ID), err
}

func DeleteTest(testID uint) error {
	if err := db.GetDBConn().Model(&models.Test{}).Delete(&models.Test{}, testID).Error; err != nil {
		logger.Error.Printf("[repository.DeleteTest] Error deleting test %v", err)
		return TranslateGormError(err)
	}

	return nil
}

func GetAllSubmissions(testID uint) ([]models.TestSubmission, error) {
	var subs []models.TestSubmission

	err := db.GetDBConn().
		Preload("Test").             // связь TestSubmission→Test
		Preload("Answers").          // связь TestSubmission→Answers
		Preload("Answers.Question"). // связь Answer→Question
		Preload("Answers.Question.Choices").
		Where("test_id = ?", testID). // связь Question→Choices
		Find(&subs).Error
	if err != nil {
		return nil, TranslateGormError(err)
	}
	return subs, nil
}

func CreateSubmission(submission models.TestSubmission) (err error) {
	if err = db.GetDBConn().Model(&models.TestSubmission{}).Create(&submission).Error; err != nil {
		logger.Error.Printf("[repository.CreateSubmission] Error creating submission %v", err)
		return TranslateGormError(err)
	}

	return nil
}
