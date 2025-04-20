package errs

import "errors"

// General Errors
var (
	ErrAddressNotFound         = errors.New("ErrAddressNotFound")
	ErrProductReviewNotFound   = errors.New("ErrProductReviewNotFound")
	ErrAccountNotFound         = errors.New("ErrAccountNotFound")
	ErrFeaturedProductNotFound = errors.New("ErrFeaturedProductNotFound")
	ErrPaymentNotFound         = errors.New("ErrPaymentNotFound")
	ErrRecordNotFound          = errors.New("ErrRecordNotFound")
	ErrProductNotFound         = errors.New("ErrProductNotFound")
	ErrOrderNotFound           = errors.New("ErrOrderNotFound")
	ErrCategoryNotFound        = errors.New("ErrCategoryNotFound")
	ErrOrderStatusNotFound     = errors.New("ErrOrderStatusNotFound")
	ErrSomethingWentWrong      = errors.New("ErrSomethingWentWrong")
	ErrNoProductFound          = errors.New("ErrNoProductFound")
	ErrStoreNotFound           = errors.New("ErrStoreNotFound")
	ErrDeleteFailed            = errors.New("ErrDeleteFailed")
	ErrCourseNotFound          = errors.New("ErrCourseNotFound")
	ErrNoVacancyFound          = errors.New("ErrNoVacancyFound")
	ErrFetchingProducts        = errors.New("ErrFetchingProducts")
	WarningNoProductsFound     = errors.New("WarningNoProductsFound")
	ErrStoreReviewNotFound     = errors.New("ErrStoreReviewNotFound")
)
