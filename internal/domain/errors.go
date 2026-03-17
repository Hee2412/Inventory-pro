package domain

import "errors"

var (
	ErrNotFound       = errors.New("the requested resource was not found")
	ErrInvalidInput   = errors.New("provided input data is invalid")
	ErrUnauthorized   = errors.New("you do not have permission to perform this action")
	ErrInternalServer = errors.New("an internal system error occurred")
	ErrDatabase       = errors.New("a database error occurred")

	ErrSessionNotFound  = errors.New("audit session does not exist")
	ErrSessionClosed    = errors.New("this audit session has already been closed")
	ErrStoreNotAssigned = errors.New("this store is not assigned to the current audit session")

	ErrStoreInactive = errors.New("this store is currently inactive")
	ErrUserNotFound  = errors.New("user account not found")
	ErrInvalidRole   = errors.New("user does not have the required role for this operation")
)
