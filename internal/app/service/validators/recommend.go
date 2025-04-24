package validators

import (
	"Ecadr/internal/app/models"
	"Ecadr/pkg/errs"
)

func ValidateRecommend(recommend models.Recommend) (err error) {
	if recommend.TargetID == 0 {
		return errs.ErrInvalidRecommendIDs
	}

	if recommend.TargetType != "course" && recommend.TargetType != "vacancy" {
		return errs.ErrInvalidType
	}

	return nil
}
