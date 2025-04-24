package service

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service/validators"
	"Ecadr/internal/repository"
)

func GetUserRecommends(userID uint) (recommends []models.Recommend, err error) {
	recommends, err = repository.GetUserRecommendation(userID)
	if err != nil {
		return nil, err
	}

	return recommends, nil
}

func GetUserRecommendsCourse(userID uint) (recommends []models.Recommend, err error) {
	recommends, err = repository.GetUserRecommendationCourse(userID)
	if err != nil {
		return nil, err
	}

	return recommends, nil
}

func GetUserRecommendsVacancy(userID uint) (recommends []models.Recommend, err error) {
	recommends, err = repository.GetUserRecommendationVacancy(userID)
	if err != nil {
		return nil, err
	}

	return recommends, nil
}

func GetUserRecommendByID(recommendID, userID uint) (recommend *models.Recommend, err error) {
	recommend, err = repository.GetUserRecommendByID(recommendID, userID)
	if err != nil {
		return nil, err
	}

	return recommend, nil
}

func CreateRecommend(recommends models.Recommend) (err error) {
	if err = validators.ValidateRecommend(recommends); err != nil {
		return err
	}

	if recommends.TargetType == "vacancy" {
		_, err = GetVacancyByID(int(recommends.TargetID))
		if err != nil {
			return err
		}
	}

	err = repository.CreateRecommend(recommends)
	if err != nil {
		return err
	}

	return nil
}
