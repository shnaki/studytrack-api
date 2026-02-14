// Package domain contains the core business logic and entities.
package domain

import (
	"errors"
	"fmt"
)

// ErrorType represents the category of a domain error.
type ErrorType int

const (
	// ErrorTypeNotFound indicates a resource was not found.
	ErrorTypeNotFound ErrorType = iota
	// ErrorTypeValidation indicates a validation error.
	ErrorTypeValidation
	// ErrorTypeConflict indicates a conflict with the current state of the resource.
	ErrorTypeConflict
)

// Error is a custom error type for domain-specific errors.
type Error struct {
	Type    ErrorType
	Message string
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// ErrNotFound creates a new Error of type ErrorTypeNotFound.
func ErrNotFound(entity string) *Error {
	return &Error{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", entity),
	}
}

// ErrValidation creates a new Error of type ErrorTypeValidation.
func ErrValidation(message string) *Error {
	return &Error{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

// ErrConflict creates a new Error of type ErrorTypeConflict.
func ErrConflict(message string) *Error {
	return &Error{
		Type:    ErrorTypeConflict,
		Message: message,
	}
}

// IsNotFound checks if the error is an Error of type ErrorTypeNotFound.
func IsNotFound(err error) bool {
	var domErr *Error
	if errors.As(err, &domErr) {
		return domErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsValidation checks if the error is an Error of type ErrorTypeValidation.
func IsValidation(err error) bool {
	var domErr *Error
	if errors.As(err, &domErr) {
		return domErr.Type == ErrorTypeValidation
	}
	return false
}

// IsConflict checks if the error is an Error of type ErrorTypeConflict.
func IsConflict(err error) bool {
	var domErr *Error
	if errors.As(err, &domErr) {
		return domErr.Type == ErrorTypeConflict
	}
	return false
}
