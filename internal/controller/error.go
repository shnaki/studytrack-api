package controller

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func toHTTPError(err error) error {
	var domErr *domain.DomainError
	if errors.As(err, &domErr) {
		switch domErr.Type {
		case domain.ErrorTypeNotFound:
			return huma.Error404NotFound(domErr.Message)
		case domain.ErrorTypeValidation:
			return huma.Error400BadRequest(domErr.Message)
		case domain.ErrorTypeConflict:
			return huma.Error409Conflict(domErr.Message)
		}
	}
	return huma.Error500InternalServerError("internal server error")
}
