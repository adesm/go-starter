package user

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
