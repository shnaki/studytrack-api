package usecase_test

import (
	"context"
	"testing"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

func setupSubjectTest() (*usecase.SubjectUsecase, *mockUserRepository, *mockSubjectRepository) {
	userRepo := newMockUserRepository()
	subjectRepo := newMockSubjectRepository()
	uc := usecase.NewSubjectUsecase(subjectRepo, userRepo)
	return uc, userRepo, subjectRepo
}

func createTestUser(repo *mockUserRepository, id, name string) {
	repo.users[id] = &domain.User{ID: id, Name: name}
}

func TestCreateSubject_Success(t *testing.T) {
	uc, userRepo, _ := setupSubjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	subject, err := uc.CreateSubject(context.Background(), "user-1", "Mathematics")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject.Name != "Mathematics" {
		t.Errorf("expected name 'Mathematics', got '%s'", subject.Name)
	}
	if subject.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", subject.UserID)
	}
	if subject.ID == "" {
		t.Error("expected ID to be generated")
	}
}

func TestCreateSubject_UserNotFound(t *testing.T) {
	uc, _, _ := setupSubjectTest()

	_, err := uc.CreateSubject(context.Background(), "nonexistent", "Math")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateSubject_EmptyName(t *testing.T) {
	uc, userRepo, _ := setupSubjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	_, err := uc.CreateSubject(context.Background(), "user-1", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestCreateSubject_DuplicateName(t *testing.T) {
	uc, userRepo, _ := setupSubjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	_, err := uc.CreateSubject(context.Background(), "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error creating first subject: %v", err)
	}

	_, err = uc.CreateSubject(context.Background(), "user-1", "Math")
	if err == nil {
		t.Fatal("expected error for duplicate subject name")
	}
	if !domain.IsConflict(err) {
		t.Errorf("expected conflict error, got: %v", err)
	}
}

func TestListSubjects_Success(t *testing.T) {
	uc, userRepo, _ := setupSubjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	_, err := uc.CreateSubject(context.Background(), "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = uc.CreateSubject(context.Background(), "user-1", "English")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	subjects, err := uc.ListSubjects(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subjects) != 2 {
		t.Errorf("expected 2 subjects, got %d", len(subjects))
	}
}

func TestListSubjects_UserNotFound(t *testing.T) {
	uc, _, _ := setupSubjectTest()

	_, err := uc.ListSubjects(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpdateSubject_Success(t *testing.T) {
	uc, userRepo, _ := setupSubjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	created, err := uc.CreateSubject(context.Background(), "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := uc.UpdateSubject(context.Background(), created.ID, "Mathematics")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Mathematics" {
		t.Errorf("expected name 'Mathematics', got '%s'", updated.Name)
	}
}

func TestUpdateSubject_NotFound(t *testing.T) {
	uc, _, _ := setupSubjectTest()

	_, err := uc.UpdateSubject(context.Background(), "nonexistent", "NewName")
	if err == nil {
		t.Fatal("expected error for nonexistent subject")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestDeleteSubject_Success(t *testing.T) {
	uc, userRepo, subjectRepo := setupSubjectTest()
	createTestUser(userRepo, "user-1", "Alice")

	created, err := uc.CreateSubject(context.Background(), "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = uc.DeleteSubject(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify deleted
	if _, ok := subjectRepo.subjects[created.ID]; ok {
		t.Error("expected subject to be deleted from repository")
	}
}

func TestDeleteSubject_NotFound(t *testing.T) {
	uc, _, _ := setupSubjectTest()

	err := uc.DeleteSubject(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent subject")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}
