package controller

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestToHTTPError_NotFound(t *testing.T) {
	domErr := domain.ErrNotFound("user")
	httpErr := toHTTPError(domErr)

	var se huma.StatusError
	if !errors.As(httpErr, &se) {
		t.Fatalf("expected huma.StatusError, got %T: %v", httpErr, httpErr)
	}
	if se.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, se.GetStatus())
	}
}

func TestToHTTPError_Validation(t *testing.T) {
	domErr := domain.ErrValidation("name is required")
	httpErr := toHTTPError(domErr)

	var se huma.StatusError
	if !errors.As(httpErr, &se) {
		t.Fatalf("expected huma.StatusError, got %T: %v", httpErr, httpErr)
	}
	if se.GetStatus() != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, se.GetStatus())
	}
}

func TestToHTTPError_Conflict(t *testing.T) {
	domErr := domain.ErrConflict("already exists")
	httpErr := toHTTPError(domErr)

	var se huma.StatusError
	if !errors.As(httpErr, &se) {
		t.Fatalf("expected huma.StatusError, got %T: %v", httpErr, httpErr)
	}
	if se.GetStatus() != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, se.GetStatus())
	}
}

func TestToHTTPError_GenericError(t *testing.T) {
	genericErr := fmt.Errorf("something went wrong")
	httpErr := toHTTPError(genericErr)

	var se huma.StatusError
	if !errors.As(httpErr, &se) {
		t.Fatalf("expected huma.StatusError, got %T: %v", httpErr, httpErr)
	}
	if se.GetStatus() != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, se.GetStatus())
	}
}

func TestToHTTPError_WrappedDomainError(t *testing.T) {
	domErr := domain.ErrNotFound("subject")
	wrapped := fmt.Errorf("wrap: %w", domErr)
	httpErr := toHTTPError(wrapped)

	var se huma.StatusError
	if !errors.As(httpErr, &se) {
		t.Fatalf("expected huma.StatusError, got %T: %v", httpErr, httpErr)
	}
	if se.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d for wrapped not-found error, got %d", http.StatusNotFound, se.GetStatus())
	}
}
