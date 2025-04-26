package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
)

func GetVacancyCriteria(vacancyID uint) (criteria []models.Criteria, err error) {
	if err = db.GetDBConn().Model(&models.Criteria{}).Where("vacancy_id = ?", vacancyID).Find(&criteria).Error; err != nil {
		logger.Error.Printf("[repository.GetVacancyCriteria] Error while getting vacancy criteria: %v", err)

		return criteria, TranslateGormError(err)
	}

	return criteria, nil
}

func GetVacancyCriteriaByID(criteriaID uint) (criteria models.Criteria, err error) {
	if err = db.GetDBConn().Model(&models.Criteria{}).Where("id = ?", criteriaID).Find(&criteria).Error; err != nil {
		logger.Error.Printf("[repository.GetVacancyCriteriaByID] Error while getting vacancy criteria: %v", err)

		return criteria, TranslateGormError(err)
	}

	return criteria, nil
}

func GetVacancyCriteriaByTitleAndVacancyID(title string, vacancyID uint) (criteria models.Criteria, err error) {
	if err = db.GetDBConn().Model(&models.Criteria{}).Where("title = ? AND vacancy_id = ?", title, vacancyID).First(&criteria).Error; err != nil {
		logger.Error.Printf("[repository.GetVacancyCriteriaByTitleAndVacancyID] Error while getting vacancy criteria: %v", err)

		return criteria, TranslateGormError(err)
	}

	return criteria, nil
}

func CreateVacancyCriteria(criteria models.Criteria) (err error) {
	if err = db.GetDBConn().Model(&models.Criteria{}).Create(&criteria).Error; err != nil {
		logger.Error.Printf("[repository.CreateVacancyCriteria] Error while creating vacancy criteria: %v", err)

		return TranslateGormError(err)
	}

	return nil
}

func UpdateVacancyCriteria(criteria models.Criteria) (err error) {
	if err = db.GetDBConn().Model(&models.Criteria{}).Where("id = ?", criteria.ID).Updates(&criteria).Error; err != nil {
		logger.Error.Printf("[repository.UpdateVacancyCriteria] Error while creating vacancy criteria: %v", err)

		return TranslateGormError(err)
	}

	return nil
}

func DeleteVacancyCriteria(criteriaID uint) (err error) {
	criteria, err := GetVacancyCriteriaByID(criteriaID)
	if err != nil {
		return TranslateGormError(err)
	}

	if err = db.GetDBConn().Model(&models.Criteria{}).Delete(&criteria).Error; err != nil {
		logger.Error.Printf("[repository.DeleteVacancyCriteria] Error while deleting vacancy criteria: %v", err)

		return TranslateGormError(err)
	}

	return nil
}
