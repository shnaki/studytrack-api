package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/repository/sqlcgen"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

type subjectRepository struct {
	q *sqlcgen.Queries
}

// NewSubjectRepository creates a new SubjectRepository implementation using PostgreSQL.
func NewSubjectRepository(pool *pgxpool.Pool) port.SubjectRepository {
	return &subjectRepository{q: sqlcgen.New(pool)}
}

func (r *subjectRepository) Create(ctx context.Context, subject *domain.Subject) error {
	err := r.q.CreateSubject(ctx, sqlcgen.CreateSubjectParams{
		ID:        toPgUUID(subject.ID),
		UserID:    toPgUUID(subject.UserID),
		Name:      subject.Name,
		CreatedAt: toPgTimestamptz(subject.CreatedAt),
		UpdatedAt: toPgTimestamptz(subject.UpdatedAt),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrConflict("subject with this name already exists for this user")
		}
		return fmt.Errorf("insert subject: %w", err)
	}
	return nil
}

func (r *subjectRepository) FindByID(ctx context.Context, id string) (*domain.Subject, error) {
	row, err := r.q.GetSubjectByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("subject")
		}
		return nil, fmt.Errorf("find subject: %w", err)
	}
	return domain.ReconstructSubject(
		fromPgUUID(row.ID),
		fromPgUUID(row.UserID),
		row.Name,
		fromPgTimestamptz(row.CreatedAt),
		fromPgTimestamptz(row.UpdatedAt),
	), nil
}

func (r *subjectRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Subject, error) {
	rows, err := r.q.ListSubjectsByUserID(ctx, toPgUUID(userID))
	if err != nil {
		return nil, fmt.Errorf("find subjects: %w", err)
	}
	subjects := make([]*domain.Subject, 0, len(rows))
	for _, row := range rows {
		subjects = append(subjects, domain.ReconstructSubject(
			fromPgUUID(row.ID),
			fromPgUUID(row.UserID),
			row.Name,
			fromPgTimestamptz(row.CreatedAt),
			fromPgTimestamptz(row.UpdatedAt),
		))
	}
	return subjects, nil
}

func (r *subjectRepository) Update(ctx context.Context, subject *domain.Subject) error {
	tag, err := r.q.UpdateSubject(ctx, sqlcgen.UpdateSubjectParams{
		Name:      subject.Name,
		UpdatedAt: toPgTimestamptz(subject.UpdatedAt),
		ID:        toPgUUID(subject.ID),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrConflict("subject with this name already exists for this user")
		}
		return fmt.Errorf("update subject: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("subject")
	}
	return nil
}

func (r *subjectRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.q.DeleteSubject(ctx, toPgUUID(id))
	if err != nil {
		return fmt.Errorf("delete subject: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("subject")
	}
	return nil
}
