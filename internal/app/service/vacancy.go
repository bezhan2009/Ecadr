package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/repository"
	"Ecadr/pkg/errs"
)

func GetAllVacancies(search string) (vacancies []models.Vacancy, err error) {
	vacancies, err = repository.GetAllVacancies(search)
	if err != nil {
		return []models.Vacancy{}, err
	}

	return vacancies, nil
}

func GetAllWorkerVacancies(workerID uint, search string) (vacancies []models.Vacancy, err error) {
	vacancies, err = repository.GetAllWorkerVacancies(workerID, search)
	if err != nil {
		return []models.Vacancy{}, err
	}

	return vacancies, nil
}

func GetVacancyByID(vacancyID int) (vacancy models.Vacancy, err error) {
	vacancy, err = repository.GetVacancyById(vacancyID)
	if err != nil {
		return models.Vacancy{}, err
	}

	return vacancy, nil
}

func CreateVacancy(vacancy models.Vacancy) (vacancyID uint, err error) {
	vacancyID, err = repository.CreateVacancy(vacancy)
	if err != nil {
		return 0, err
	}

	return vacancyID, nil
}

func UpdateVacancy(userID, vacancyID int, vacancyReq models.Vacancy) (err error) {
	vacancy, err := GetVacancyByID(vacancyID)
	if err != nil {
		return err
	}

	if vacancy.WorkerID != userID {
		return errs.ErrRecordNotFound
	}

	if vacancyReq.Title != "" {
		vacancy.Title = vacancyReq.Title
	}

	if vacancyReq.Description != "" {
		vacancy.Description = vacancyReq.Description
	}

	if vacancyReq.Contact != "" {
		vacancy.Contact = vacancyReq.Contact
	}

	err = repository.UpdateVacancy(vacancy)
	if err != nil {
		return err
	}

	return nil
}

func DeleteVacancy(userID, vacancyID int) (err error) {
	vacancy, err := GetVacancyByID(vacancyID)
	if err != nil {
		return err
	}

	if vacancy.WorkerID != userID {
		return errs.ErrRecordNotFound
	}

	err = repository.DeleteVacancy(vacancy)
	if err != nil {
		return err
	}

	return nil
}
