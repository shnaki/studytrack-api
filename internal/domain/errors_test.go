package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestDomainError_Error(t *testing.T) {
	err := domain.ErrNotFound("user")
	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found', got '%s'", err.Error())
	}

	err = domain.ErrValidation("name is required")
	if err.Error() != "name is required" {
		t.Errorf("expected 'name is required', got '%s'", err.Error())
	}

	err = domain.ErrConflict("already exists")
	if err.Error() != "already exists" {
		t.Errorf("expected 'already exists', got '%s'", err.Error())
	}
}

func TestErrNotFound(t *testing.T) {
	err := domain.ErrNotFound("subject")
	if err.Error() != "subject not found" {
		t.Errorf("expected 'subject not found', got '%s'", err.Error())
	}
	if !domain.IsNotFound(err) {
		t.Error("expected IsNotFound to return true for ErrNotFound")
	}
	if domain.IsValidation(err) {
		t.Error("expected IsValidation to return false for ErrNotFound")
	}
	if domain.IsConflict(err) {
		t.Error("expected IsConflict to return false for ErrNotFound")
	}
}

func TestErrValidation(t *testing.T) {
	err := domain.ErrValidation("invalid input")
	if err.Error() != "invalid input" {
		t.Errorf("expected 'invalid input', got '%s'", err.Error())
	}
	if !domain.IsValidation(err) {
		t.Error("expected IsValidation to return true for ErrValidation")
	}
	if domain.IsNotFound(err) {
		t.Error("expected IsNotFound to return false for ErrValidation")
	}
	if domain.IsConflict(err) {
		t.Error("expected IsConflict to return false for ErrValidation")
	}
}

func TestErrConflict(t *testing.T) {
	err := domain.ErrConflict("duplicate entry")
	if err.Error() != "duplicate entry" {
		t.Errorf("expected 'duplicate entry', got '%s'", err.Error())
	}
	if !domain.IsConflict(err) {
		t.Error("expected IsConflict to return true for ErrConflict")
	}
	if domain.IsNotFound(err) {
		t.Error("expected IsNotFound to return false for ErrConflict")
	}
	if domain.IsValidation(err) {
		t.Error("expected IsValidation to return false for ErrConflict")
	}
}

func TestIsNotFound_WithNonDomainError(t *testing.T) {
	err := errors.New("some other error")
	if domain.IsNotFound(err) {
		t.Error("expected IsNotFound to return false for non-Error")
	}
}

func TestIsValidation_WithNonDomainError(t *testing.T) {
	err := errors.New("some other error")
	if domain.IsValidation(err) {
		t.Error("expected IsValidation to return false for non-Error")
	}
}

func TestIsConflict_WithNonDomainError(t *testing.T) {
	err := errors.New("some other error")
	if domain.IsConflict(err) {
		t.Error("expected IsConflict to return false for non-Error")
	}
}

func TestErrorsAs_WrappedError(t *testing.T) {
	err := domain.ErrNotFound("user")
	// fmt.Errorf("%w", err) 形式でラップしないと As は機能しない
	wrappedErr := fmt.Errorf("wrapped: %w", err)

	if !domain.IsNotFound(wrappedErr) {
		t.Error("expected IsNotFound to return true for wrapped ErrNotFound")
	}
}
