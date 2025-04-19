package validators

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
)

func ValidateCourse(course models.Course) error {
	if len(course.Title) < 5 {
		return errs.ErrInvalidTitle
	}

	if len(course.Description) < 10 {
		return errs.ErrInvalidDescription
	}

	// Проверка, что даты не нулевые
	if course.StartDate.IsZero() || course.EndDate.IsZero() {
		return errs.ErrInvalidDate
	}

	// Проверка, что дата конца после даты начала
	if course.EndDate.Before(course.StartDate) || course.EndDate.Equal(course.StartDate) {
		return errs.ErrEndDateBeforeStartDate
	}

	return nil
}
