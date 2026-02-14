package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

func setupGoalTest() (*usecase.GoalUsecase, *mockUserRepository, *mockSubjectRepository, *mockGoalRepository) {
	userRepo := newMockUserRepository()
	subjectRepo := newMockSubjectRepository()
	goalRepo := newMockGoalRepository()
	uc := usecase.NewGoalUsecase(goalRepo, userRepo, subjectRepo)
	return uc, userRepo, subjectRepo, goalRepo
}

func TestUpsertGoal_CreateNew(t *testing.T) {
	uc, userRepo, subjectRepo, goalRepo := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")
	subjectRepo.subjects["sub-1"] = &domain.Subject{ID: "sub-1", UserID: "user-1", Name: "Math"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	goal, err := uc.UpsertGoal(context.Background(), "user-1", "sub-1", 300, startDate, &endDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", goal.UserID)
	}
	if goal.SubjectID != "sub-1" {
		t.Errorf("expected SubjectID 'sub-1', got '%s'", goal.SubjectID)
	}
	if goal.TargetMinutesPerWeek != 300 {
		t.Errorf("expected TargetMinutesPerWeek 300, got %d", goal.TargetMinutesPerWeek)
	}
	if goal.ID == "" {
		t.Error("expected ID to be generated")
	}
	if len(goalRepo.goals) != 1 {
		t.Errorf("expected 1 goal in repo, got %d", len(goalRepo.goals))
	}
}

func TestUpsertGoal_UserNotFound(t *testing.T) {
	uc, _, subjectRepo, _ := setupGoalTest()
	subjectRepo.subjects["sub-1"] = &domain.Subject{ID: "sub-1", UserID: "user-1", Name: "Math"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := uc.UpsertGoal(context.Background(), "nonexistent", "sub-1", 300, startDate, nil)
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpsertGoal_SubjectNotFound(t *testing.T) {
	uc, userRepo, _, _ := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := uc.UpsertGoal(context.Background(), "user-1", "nonexistent", 300, startDate, nil)
	if err == nil {
		t.Fatal("expected error for nonexistent subject")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpsertGoal_SubjectBelongsToDifferentUser(t *testing.T) {
	uc, userRepo, subjectRepo, _ := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestUser(userRepo, "user-2", "Bob")
	subjectRepo.subjects["sub-1"] = &domain.Subject{ID: "sub-1", UserID: "user-2", Name: "Math"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := uc.UpsertGoal(context.Background(), "user-1", "sub-1", 300, startDate, nil)
	if err == nil {
		t.Fatal("expected error when subject belongs to different user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpsertGoal_InvalidTarget(t *testing.T) {
	uc, userRepo, subjectRepo, _ := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")
	subjectRepo.subjects["sub-1"] = &domain.Subject{ID: "sub-1", UserID: "user-1", Name: "Math"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Zero target
	_, err := uc.UpsertGoal(context.Background(), "user-1", "sub-1", 0, startDate, nil)
	if err == nil {
		t.Fatal("expected error for zero target")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error for zero target, got: %v", err)
	}

	// Negative target
	_, err = uc.UpsertGoal(context.Background(), "user-1", "sub-1", -10, startDate, nil)
	if err == nil {
		t.Fatal("expected error for negative target")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error for negative target, got: %v", err)
	}
}

func TestListGoals_Success(t *testing.T) {
	uc, userRepo, subjectRepo, _ := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")
	subjectRepo.subjects["sub-1"] = &domain.Subject{ID: "sub-1", UserID: "user-1", Name: "Math"}
	subjectRepo.subjects["sub-2"] = &domain.Subject{ID: "sub-2", UserID: "user-1", Name: "English"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := uc.UpsertGoal(context.Background(), "user-1", "sub-1", 300, startDate, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = uc.UpsertGoal(context.Background(), "user-1", "sub-2", 120, startDate, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	goals, err := uc.ListGoals(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(goals) != 2 {
		t.Errorf("expected 2 goals, got %d", len(goals))
	}
}

func TestListGoals_UserNotFound(t *testing.T) {
	uc, _, _, _ := setupGoalTest()

	_, err := uc.ListGoals(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}
