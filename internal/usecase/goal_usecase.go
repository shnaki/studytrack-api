// Package usecase contains the application logic.
package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

// GoalUsecase provides methods for managing study goals.
type GoalUsecase struct {
	goalRepo    port.GoalRepository
	userRepo    port.UserRepository
	projectRepo port.ProjectRepository
}

// NewGoalUsecase creates a new GoalUsecase.
func NewGoalUsecase(
	goalRepo port.GoalRepository,
	userRepo port.UserRepository,
	projectRepo port.ProjectRepository,
) *GoalUsecase {
	return &GoalUsecase{
		goalRepo:    goalRepo,
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

// UpsertGoal creates or updates a goal for a project.
func (u *GoalUsecase) UpsertGoal(ctx context.Context, userID, projectID string, targetMinutesPerWeek int, startDate time.Time, endDate *time.Time) (*domain.Goal, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	project, err := u.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, domain.ErrNotFound("project")
	}

	id := uuid.New().String()
	goal, err := domain.NewGoal(id, userID, projectID, targetMinutesPerWeek, startDate, endDate)
	if err != nil {
		return nil, err
	}
	if err := u.goalRepo.Upsert(ctx, goal); err != nil {
		return nil, err
	}
	return goal, nil
}

// ListGoals returns all goals for a user.
func (u *GoalUsecase) ListGoals(ctx context.Context, userID string) ([]*domain.Goal, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	return u.goalRepo.FindByUserID(ctx, userID)
}
