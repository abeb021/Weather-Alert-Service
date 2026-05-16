package errors

import "errors"

var (
	ErrEmptyEmail    = errors.New("email is empty")
	ErrInvalidEmail  = errors.New("email is invalid")
	ErrEmptySubject  = errors.New("subject is empty")
	ErrEmptyMessage  = errors.New("message is empty")
	ErrEmailSendFail = errors.New("failed to send email")
)
