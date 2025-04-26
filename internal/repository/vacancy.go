package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/logger"
)

func GetAllVacancies(search string) (vacancies []models.Vacancy, err error) {
	query := db.GetDBConn().Preload("Company").Preload("Criteria")

	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", likeSearch, likeSearch)
	}

	if err = query.Find(&vacancies).Error; err != nil {
		logger.Error.Printf("[repository.GetAllVacancies] Err while getting vacancies: %v", err)

		return nil, TranslateGormError(err)
	}

	return vacancies, nil
}

func GetAllWorkerVacancies(workerID uint, search string) (vacancies []models.Vacancy, err error) {
	query := db.GetDBConn().Preload("Company").Preload("Criteria")

	if search != "" {
		// Пример: ищем по title и description
		likeSearch := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", likeSearch, likeSearch)
	}

	if err = query.Where("worker_id = ?", workerID).Find(&vacancies).Error; err != nil {
		logger.Error.Printf("[repository.GetAllWorkerVacancies] Err while getting worker vacancies: %v", err)
		return nil, TranslateGormError(err)
	}

	return vacancies, nil
}

func GetVacancyById(vacancyId int) (vacancy models.Vacancy, err error) {
	if err = db.GetDBConn().Preload("Company").Preload("Criteria").Where("id = ?", vacancyId).First(&vacancy).Error; err != nil {
		logger.Error.Printf("[repository.GetVacancyById] Err while getting vacancy: %v", err)

		return models.Vacancy{}, TranslateGormError(err)
	}

	return vacancy, nil
}

func CreateVacancy(vacancy models.Vacancy) (uint, error) {
	if err := db.GetDBConn().Create(&vacancy).Error; err != nil {
		logger.Error.Printf("[repository.CreateVacancy] Err while creating vacancy: %v", err)

		return vacancy.ID, TranslateGormError(err)
	}

	return vacancy.ID, nil
}

func UpdateVacancy(vacancy models.Vacancy) (err error) {
	if err = db.GetDBConn().Updates(&vacancy).Error; err != nil {
		logger.Error.Printf("[repository.UpdateVacancy] Err while updating vacancy: %v", err)

		return TranslateGormError(err)
	}

	return nil
}

func DeleteVacancy(vacancy models.Vacancy) (err error) {
	if err = db.GetDBConn().Delete(&vacancy).Error; err != nil {
		logger.Error.Printf("[repository.DeleteVacancy] Err while deleting vacancy: %v", err)

		return TranslateGormError(err)
	}

	return nil
}
