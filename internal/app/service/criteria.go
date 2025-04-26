package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/repository"
	"Ecadr/pkg/errs"
)

func GetVacancyCriteria(vacancyID uint) (criteria []models.Criteria, err error) {
	criteria, err = repository.GetVacancyCriteria(vacancyID)
	if err != nil {
		return nil, err
	}

	return criteria, nil
}

func GetVacancyCriteriaByID(criteriaID uint) (criteria models.Criteria, err error) {
	criteria, err = repository.GetVacancyCriteriaByID(criteriaID)
	if err != nil {
		return models.Criteria{}, err
	}

	return criteria, nil
}

func CreateCriteria(criteria models.Criteria) (err error) {
	if criteria.Title == "" {
		return errs.ErrInvalidTitle
	}

	_, err = repository.GetVacancyById(int(criteria.VacancyID))
	if err != nil {
		return err
	}

	_, err = repository.GetVacancyCriteriaByTitleAndVacancyID(criteria.Title, criteria.VacancyID)
	if err == nil {
		return errs.ErrAlreadyExists
	}

	err = repository.CreateVacancyCriteria(criteria)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCriteria(criteria models.Criteria) (err error) {
	_, err = repository.GetVacancyCriteriaByID(criteria.ID)
	if err != nil {
		return err
	}

	err = repository.UpdateVacancyCriteria(criteria)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCriteria(criteriaID uint) (err error) {
	err = repository.DeleteVacancyCriteria(criteriaID)
	if err != nil {
		return err
	}

	return nil
}
