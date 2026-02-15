package usecase_test

import (
	"context"
	"testing"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

func setupProjectTest() (*usecase.ProjectUsecase, *mockUserRepository, *mockProjectRepository) {
	userRepo := newMockUserRepository()
	projectRepo := newMockProjectRepository()
	uc := usecase.NewProjectUsecase(projectRepo, userRepo)
	return uc, userRepo, projectRepo
}

func createTestUser(repo *mockUserRepository, id, name string) {
	repo.users[id] = &domain.User{ID: id, Name: name}
}

func TestCreateProject_Success(t *testing.T) {
	uc, userRepo, _ := setupProjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	project, err := uc.CreateProject(context.Background(), "user-1", "Mathematics")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.Name != "Mathematics" {
		t.Errorf("expected name 'Mathematics', got '%s'", project.Name)
	}
	if project.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", project.UserID)
	}
	if project.ID == "" {
		t.Error("expected ID to be generated")
	}
}

func TestCreateProject_UserNotFound(t *testing.T) {
	uc, _, _ := setupProjectTest()

	_, err := uc.CreateProject(context.Background(), "nonexistent", "Math")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateProject_EmptyName(t *testing.T) {
	uc, userRepo, _ := setupProjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	_, err := uc.CreateProject(context.Background(), "user-1", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestCreateProject_DuplicateName(t *testing.T) {
	uc, userRepo, _ := setupProjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	_, err := uc.CreateProject(context.Background(), "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error creating first project: %v", err)
	}

	_, err = uc.CreateProject(context.Background(), "user-1", "Math")
	if err == nil {
		t.Fatal("expected error for duplicate project name")
	}
	if !domain.IsConflict(err) {
		t.Errorf("expected conflict error, got: %v", err)
	}
}

func TestListProjects_Success(t *testing.T) {
	uc, userRepo, _ := setupProjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	_, err := uc.CreateProject(context.Background(), "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = uc.CreateProject(context.Background(), "user-1", "English")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	projects, err := uc.ListProjects(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
}

func TestListProjects_UserNotFound(t *testing.T) {
	uc, _, _ := setupProjectTest()

	_, err := uc.ListProjects(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpdateProject_Success(t *testing.T) {
	uc, userRepo, _ := setupProjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	created, err := uc.CreateProject(context.Background(), "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := uc.UpdateProject(context.Background(), created.ID, "Mathematics")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Mathematics" {
		t.Errorf("expected name 'Mathematics', got '%s'", updated.Name)
	}
}

func TestUpdateProject_NotFound(t *testing.T) {
	uc, _, _ := setupProjectTest()

	_, err := uc.UpdateProject(context.Background(), "nonexistent", "NewName")
	if err == nil {
		t.Fatal("expected error for nonexistent project")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestDeleteProject_Success(t *testing.T) {
	uc, userRepo, projectRepo := setupProjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	created, err := uc.CreateProject(context.Background(), "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = uc.DeleteProject(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify deleted
	if _, ok := projectRepo.projects[created.ID]; ok {
		t.Error("expected project to be deleted from repository")
	}
}

func TestDeleteProject_NotFound(t *testing.T) {
	uc, _, _ := setupProjectTest()

	err := uc.DeleteProject(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent project")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}
