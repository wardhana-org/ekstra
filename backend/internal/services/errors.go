package services

import "errors"

var (
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
