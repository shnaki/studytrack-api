package domain_test

import (
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestNewUser_Valid(t *testing.T) {
	user, err := domain.NewUser("test-id", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", user.ID)
	}
	if user.Name != "Alice" {
		t.Errorf("expected Name 'Alice', got '%s'", user.Name)
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNewUser_EmptyName(t *testing.T) {
	_, err := domain.NewUser("test-id", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewUser_NameTooLong(t *testing.T) {
	longName := make([]byte, 101)
	for i := range longName {
		longName[i] = 'a'
	}
	_, err := domain.NewUser("test-id", string(longName))
	if err == nil {
		t.Fatal("expected error for long name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestReconstructUser(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)

	user := domain.ReconstructUser("user-1", "Alice", createdAt, updatedAt)

	if user.ID != "user-1" {
		t.Errorf("expected ID 'user-1', got '%s'", user.ID)
	}
	if user.Name != "Alice" {
		t.Errorf("expected Name 'Alice', got '%s'", user.Name)
	}
	if !user.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, user.CreatedAt)
	}
	if !user.UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, user.UpdatedAt)
	}
}
