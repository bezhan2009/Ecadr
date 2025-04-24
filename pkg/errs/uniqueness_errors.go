package errs

import "errors"

// Uniqueness Errors
var (
	ErrUsernameUniquenessFailed = errors.New("ErrUsernameUniquenessFailed")
	ErrEmailUniquenessFailed    = errors.New("ErrEmailUniquenessFailed")
)
