package services

import "errors"

var (
	ErrEmailRequired            = errors.New("email is required")
	ErrPasswordRequired         = errors.New("password is required")
	ErrPasswordTooShort         = errors.New("password is too short")
	ErrPasswordTooLong          = errors.New("password is too long")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrUsernameRequired         = errors.New("username is required")
	ErrRegisterExistingEmail    = errors.New("email already registered")
	ErrRegisterExistingUsername = errors.New("username already used")
)
