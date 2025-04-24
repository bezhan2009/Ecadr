package errs

import "errors"

// Authentication Errors
var (
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrIncorrectUsernameOrPassword = errors.New("ErrIncorrectUsernameOrPassword")
	ErrPasswordIsEmpty             = errors.New("ErrPasswordIsEmpty")
	ErrUsernameIsEmpty             = errors.New("ErrUsernameIsEmpty")
	ErrEmailIsEmpty                = errors.New("ErrEmailIsEmpty")
	ErrPermissionDenied            = errors.New("ErrPermissionDenied")
	ErrUnauthorized                = errors.New("ErrUnauthorized")
	ErrUserNotFound                = errors.New("ErrUserNotFound")
	ErrEmailIsRequired             = errors.New("email is required")
	ErrUsernameIsRequired          = errors.New("username is required")
	ErrFirstNameIsRequired         = errors.New("first name is required")
	ErrLastNameIsRequired          = errors.New("last name is required")
	ErrPasswordIsRequired          = errors.New("password is required")
	ErrWrongRoleID                 = errors.New("wrong role")
)
