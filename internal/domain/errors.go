package domain

import "errors"

var (
	// General
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")

	// User
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")

	// Product
	ErrProductNotFound = errors.New("product not found")

	// Session
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionClosed   = errors.New("session is closed")

	// Order
	ErrOrderNotFound = errors.New("order not found")
	ErrEmptyOrder    = errors.New("cannot submit empty order")
	ErrOrderClosed   = errors.New("cannot update closed order")

	// Audit
	ErrAuditNotFound = errors.New("audit not found")
	ErrAuditExpired  = errors.New("audit session expired")

	// Database
	ErrDatabase       = errors.New("database error")
	ErrInternalServer = errors.New("internal server error")

	// Transfer
	ErrTransferNotFound = errors.New("transfer not found")
)
