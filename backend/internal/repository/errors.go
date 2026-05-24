package repository

import "errors"

var (
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrDuplicateUsername  = errors.New("username already exists")
	ErrNotFound           = errors.New("record not found")
	ErrRefreshTokenRace   = errors.New("refresh token race")
	ErrRefreshTokenReused = errors.New("refresh token reused")
)
