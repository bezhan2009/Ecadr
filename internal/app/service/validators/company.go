package validators

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
)

func ValidateCompany(company models.Company) (err error) {
	if company.Title == "" {
		return errs.ErrInvalidTitle
	}

	if company.Description == "" {
		return errs.ErrInvalidDescription
	}

	return nil
}
