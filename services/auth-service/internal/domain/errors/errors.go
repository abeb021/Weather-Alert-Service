package errors

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrEmailAlreadyExists  = errors.New("email already registered")

	ErrRefreshTokenNotFound = errors.New("refresh token not found")

	ErrUserCreate   = errors.New("failed to create user")
	ErrUserNotFound = errors.New("failed to find user")
)
