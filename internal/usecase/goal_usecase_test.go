package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

func setupGoalTest() (*usecase.GoalUsecase, *mockUserRepository, *mockProjectRepository, *mockGoalRepository) {
	userRepo := newMockUserRepository()
	projectRepo := newMockProjectRepository()
	goalRepo := newMockGoalRepository()
	uc := usecase.NewGoalUsecase(goalRepo, userRepo, projectRepo)
	return uc, userRepo, projectRepo, goalRepo
}

func TestUpsertGoal_CreateNew(t *testing.T) {
	uc, userRepo, projectRepo, goalRepo := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	goal, err := uc.UpsertGoal(context.Background(), "user-1", "proj-1", 300, startDate, &endDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", goal.UserID)
	}
	if goal.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID 'proj-1', got '%s'", goal.ProjectID)
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
	uc, _, projectRepo, _ := setupGoalTest()
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := uc.UpsertGoal(context.Background(), "nonexistent", "proj-1", 300, startDate, nil)
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpsertGoal_ProjectNotFound(t *testing.T) {
	uc, userRepo, _, _ := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := uc.UpsertGoal(context.Background(), "user-1", "nonexistent", 300, startDate, nil)
	if err == nil {
		t.Fatal("expected error for nonexistent project")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpsertGoal_ProjectBelongsToDifferentUser(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestUser(userRepo, "user-2", "Bob")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-2", Name: "Math"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := uc.UpsertGoal(context.Background(), "user-1", "proj-1", 300, startDate, nil)
	if err == nil {
		t.Fatal("expected error when project belongs to different user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpsertGoal_InvalidTarget(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Zero target
	_, err := uc.UpsertGoal(context.Background(), "user-1", "proj-1", 0, startDate, nil)
	if err == nil {
		t.Fatal("expected error for zero target")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error for zero target, got: %v", err)
	}

	// Negative target
	_, err = uc.UpsertGoal(context.Background(), "user-1", "proj-1", -10, startDate, nil)
	if err == nil {
		t.Fatal("expected error for negative target")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error for negative target, got: %v", err)
	}
}

func TestListGoals_Success(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupGoalTest()
	createTestUser(userRepo, "user-1", "Alice")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}
	projectRepo.projects["proj-2"] = &domain.Project{ID: "proj-2", UserID: "user-1", Name: "English"}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := uc.UpsertGoal(context.Background(), "user-1", "proj-1", 300, startDate, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = uc.UpsertGoal(context.Background(), "user-1", "proj-2", 120, startDate, nil)
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
