package port

import (
	"context"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
}

type SubjectRepository interface {
	Create(ctx context.Context, subject *domain.Subject) error
	FindByID(ctx context.Context, id string) (*domain.Subject, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Subject, error)
	Update(ctx context.Context, subject *domain.Subject) error
	Delete(ctx context.Context, id string) error
}

type StudyLogFilter struct {
	From      *time.Time
	To        *time.Time
	SubjectID *string
}

type StudyLogRepository interface {
	Create(ctx context.Context, log *domain.StudyLog) error
	FindByID(ctx context.Context, id string) (*domain.StudyLog, error)
	FindByUserID(ctx context.Context, userID string, filter StudyLogFilter) ([]*domain.StudyLog, error)
	Delete(ctx context.Context, id string) error
}

type GoalRepository interface {
	Upsert(ctx context.Context, goal *domain.Goal) error
	FindByUserID(ctx context.Context, userID string) ([]*domain.Goal, error)
}
