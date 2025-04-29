package errs

import "errors"

// General Errors
var (
	ErrRecordNotFound        = errors.New("ErrRecordNotFound")
	ErrEmptyMessageText      = errors.New("ErrEmptyMessageText")
	ErrNoAIRecommends        = errors.New("ErrNoAIRecommends")
	ErrAlreadyExists         = errors.New("ErrAlreadyExists")
	ErrRoleIsRequired        = errors.New("ErrRoleIsRequired")
	ErrSomethingWentWrong    = errors.New("ErrSomethingWentWrong")
	ErrGeminiIsNotWorking    = errors.New("ErrGeminiIsNotWorking")
	ErrDeleteFailed          = errors.New("ErrDeleteFailed")
	ErrNoCourseFound         = errors.New("ErrNoCourseFound")
	ErrNoVacancyFound        = errors.New("ErrNoVacancyFound")
	ErrNoUsersStatisticFound = errors.New("ErrNoUsersStatisticFound")
)
