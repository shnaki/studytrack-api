package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

type GoalUsecase struct {
	goalRepo    port.GoalRepository
	userRepo    port.UserRepository
	subjectRepo port.SubjectRepository
}

func NewGoalUsecase(
	goalRepo port.GoalRepository,
	userRepo port.UserRepository,
	subjectRepo port.SubjectRepository,
) *GoalUsecase {
	return &GoalUsecase{
		goalRepo:    goalRepo,
		userRepo:    userRepo,
		subjectRepo: subjectRepo,
	}
}

func (u *GoalUsecase) UpsertGoal(ctx context.Context, userID, subjectID string, targetMinutesPerWeek int, startDate time.Time, endDate *time.Time) (*domain.Goal, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	subject, err := u.subjectRepo.FindByID(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	if subject.UserID != userID {
		return nil, domain.ErrNotFound("subject")
	}

	id := uuid.New().String()
	goal, err := domain.NewGoal(id, userID, subjectID, targetMinutesPerWeek, startDate, endDate)
	if err != nil {
		return nil, err
	}
	if err := u.goalRepo.Upsert(ctx, goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (u *GoalUsecase) ListGoals(ctx context.Context, userID string) ([]*domain.Goal, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	return u.goalRepo.FindByUserID(ctx, userID)
}
