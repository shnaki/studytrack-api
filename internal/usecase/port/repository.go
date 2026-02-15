package port

import (
	"context"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
}

// ProjectRepository defines the interface for project persistence.
type ProjectRepository interface {
	Create(ctx context.Context, project *domain.Project) error
	FindByID(ctx context.Context, id string) (*domain.Project, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Project, error)
	Update(ctx context.Context, project *domain.Project) error
	Delete(ctx context.Context, id string) error
}

// StudyLogFilter defines filters for study log queries.
type StudyLogFilter struct {
	From      *time.Time
	To        *time.Time
	ProjectID *string
}

// StudyLogRepository defines the interface for study log persistence.
type StudyLogRepository interface {
	Create(ctx context.Context, log *domain.StudyLog) error
	FindByID(ctx context.Context, id string) (*domain.StudyLog, error)
	FindByUserID(ctx context.Context, userID string, filter StudyLogFilter) ([]*domain.StudyLog, error)
	Delete(ctx context.Context, id string) error
}

// GoalRepository defines the interface for goal persistence.
type GoalRepository interface {
	Upsert(ctx context.Context, goal *domain.Goal) error
	FindByUserID(ctx context.Context, userID string) ([]*domain.Goal, error)
}

// NoteRepository defines the interface for note persistence.
type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) error
	FindByID(ctx context.Context, id string) (*domain.Note, error)
	FindByProjectID(ctx context.Context, projectID string) ([]*domain.Note, error)
	Update(ctx context.Context, note *domain.Note) error
	Delete(ctx context.Context, id string) error
}
