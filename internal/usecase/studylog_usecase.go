package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

// StudyLogUsecase provides methods for managing study logs.
type StudyLogUsecase struct {
	studyLogRepo port.StudyLogRepository
	userRepo     port.UserRepository
	projectRepo  port.ProjectRepository
}

// NewStudyLogUsecase creates a new StudyLogUsecase.
func NewStudyLogUsecase(
	studyLogRepo port.StudyLogRepository,
	userRepo port.UserRepository,
	projectRepo port.ProjectRepository,
) *StudyLogUsecase {
	return &StudyLogUsecase{
		studyLogRepo: studyLogRepo,
		userRepo:     userRepo,
		projectRepo:  projectRepo,
	}
}

// CreateStudyLog creates a new study log.
func (u *StudyLogUsecase) CreateStudyLog(ctx context.Context, userID, projectID string, studiedAt time.Time, minutes int, note string) (*domain.StudyLog, error) {
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
	log, err := domain.NewStudyLog(id, userID, projectID, studiedAt, minutes, note)
	if err != nil {
		return nil, err
	}
	if err := u.studyLogRepo.Create(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// ListStudyLogs returns all study logs for a user based on filters.
func (u *StudyLogUsecase) ListStudyLogs(ctx context.Context, userID string, filter port.StudyLogFilter) ([]*domain.StudyLog, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	return u.studyLogRepo.FindByUserID(ctx, userID, filter)
}

// DeleteStudyLog deletes a study log by ID.
func (u *StudyLogUsecase) DeleteStudyLog(ctx context.Context, id string) error {
	if _, err := u.studyLogRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return u.studyLogRepo.Delete(ctx, id)
}
