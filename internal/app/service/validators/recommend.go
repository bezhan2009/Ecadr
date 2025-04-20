package validators

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
)

func ValidateRecommendCourse(recommend models.Recommend) (err error) {
	if recommend.CourseID == 0 {
		return errs.ErrInvalidRecommendIDs
	}

	return nil
}

func ValidateRecommendVacancy(recommend models.Recommend) (err error) {
	if recommend.VacancyID == 0 {
		return errs.ErrInvalidRecommendIDs
	}

	return nil
}
