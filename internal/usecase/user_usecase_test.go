package usecase_test

import (
	"context"
	"testing"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

func TestCreateUser_Success(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	user, err := uc.CreateUser(context.Background(), "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected name 'Alice', got '%s'", user.Name)
	}
	if user.ID == "" {
		t.Error("expected ID to be generated")
	}

	// Verify stored
	stored, err := repo.FindByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("user not found in repo: %v", err)
	}
	if stored.Name != "Alice" {
		t.Errorf("stored name mismatch")
	}
}

func TestCreateUser_EmptyName(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	_, err := uc.CreateUser(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestGetUser_Found(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	created, _ := uc.CreateUser(context.Background(), "Bob")
	found, err := uc.GetUser(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != "Bob" {
		t.Errorf("expected 'Bob', got '%s'", found.Name)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	_, err := uc.GetUser(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}
