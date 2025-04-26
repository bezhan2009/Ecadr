package repository

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/db"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
)

func GetUserRecommendation(userID uint) (recommends []models.Recommend, err error) {
	if err = db.GetDBConn().Model(&models.Recommend{}).Where("user_id = ? AND vacancy_id != 0", userID).First(&recommends).Error; err != nil {
		logger.Error.Printf("[repository.GetUserRecommendation] Error while getting user recommends: %v", err)

		return recommends, TranslateGormError(err)
	}

	return recommends, nil
}

func GetUserRecommendationCourse(userID uint) (recommends []models.Recommend, err error) {
	if err = db.GetDBConn().Model(&models.Recommend{}).Where("user_id = ? AND target_type = ?", userID, "course").First(&recommends).Error; err != nil {
		logger.Error.Printf("[repository.GetUserRecommendation] Error while getting user recommends: %v", err)

		return recommends, TranslateGormError(err)
	}

	return recommends, nil
}

func GetUserRecommendationVacancy(userID uint) (recommends []models.Recommend, err error) {
	if err = db.GetDBConn().Model(&models.Recommend{}).Where("user_id = ? AND target_type = ?", userID, "vacancy").First(&recommends).Error; err != nil {
		logger.Error.Printf("[repository.GetUserRecommendation] Error while getting user recommends: %v", err)

		return recommends, TranslateGormError(err)
	}

	return recommends, nil
}

func GetUserRecommendByID(recommendID, userID uint) (recommend *models.Recommend, err error) {
	if err = db.GetDBConn().Model(&models.Recommend{}).Where("id = ? AND user_id = ?", recommendID, userID).First(&recommend).Error; err != nil {
		logger.Error.Printf("[repository.GetUserRecommendByID] Error while getting user recommends: %v", err)

		return nil, TranslateGormError(err)
	}

	return recommend, nil
}

func GetRecommendByTargetAndUserID(recommend models.Recommend) (recommendMd models.Recommend, err error) {
	if err = db.GetDBConn().Model(&models.Recommend{}).Where("target_id = ? AND target_type = ? AND user_id = ?", recommend.TargetID, recommend.TargetType, recommend.UserID).First(&recommendMd).Error; err != nil {
		logger.Error.Printf("[repository.GetRecommendByTargetAndUserID] Error while getting recommend: %v", err)

		return recommendMd, TranslateGormError(err)
	}

	return recommendMd, nil
}

func CreateRecommend(recommends models.Recommend) (err error) {
	_, err = GetRecommendByTargetAndUserID(recommends)
	if err == nil {
		return errs.ErrAlreadyExists
	}
	if err = db.GetDBConn().Model(&models.Recommend{}).Create(&recommends).Error; err != nil {
		logger.Error.Printf("[repository.CreateRecommend] Error while creating recommends: %v", err)

		return TranslateGormError(err)
	}

	return nil
}
