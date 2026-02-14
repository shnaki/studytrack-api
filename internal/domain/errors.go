package domain

import "fmt"

type ErrorType int

const (
	ErrorTypeNotFound ErrorType = iota
	ErrorTypeValidation
	ErrorTypeConflict
)

type DomainError struct {
	Type    ErrorType
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func ErrNotFound(entity string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", entity),
	}
}

func ErrValidation(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

func ErrConflict(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeConflict,
		Message: message,
	}
}

func IsNotFound(err error) bool {
	if domErr, ok := err.(*DomainError); ok {
		return domErr.Type == ErrorTypeNotFound
	}
	return false
}

func IsValidation(err error) bool {
	if domErr, ok := err.(*DomainError); ok {
		return domErr.Type == ErrorTypeValidation
	}
	return false
}

func IsConflict(err error) bool {
	if domErr, ok := err.(*DomainError); ok {
		return domErr.Type == ErrorTypeConflict
	}
	return false
}
