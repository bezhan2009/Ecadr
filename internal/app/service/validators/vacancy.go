package validators

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
)

func ValidateVacancy(vacancy models.Vacancy) (err error) {
	if len(vacancy.Title) < 5 {
		return errs.ErrInvalidTitle
	}

	if len(vacancy.Description) < 10 {
		return errs.ErrInvalidDescription
	}

	return nil
}
