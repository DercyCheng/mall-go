package repository

import "errors"

// Common repository errors
var (
	// ErrNotFound indicates that a requested resource was not found
	ErrNotFound = errors.New("resource not found")
	
	// ErrConflict indicates that a resource already exists
	ErrConflict = errors.New("resource already exists")
	
	// ErrInvalidData indicates that invalid data was provided
	ErrInvalidData = errors.New("invalid data provided")
)
