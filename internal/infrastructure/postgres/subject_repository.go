package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/application/port"
	"github.com/shnaki/studytrack-api/internal/domain"
)

type subjectRepository struct {
	pool *pgxpool.Pool
}

func NewSubjectRepository(pool *pgxpool.Pool) port.SubjectRepository {
	return &subjectRepository{pool: pool}
}

func (r *subjectRepository) Create(ctx context.Context, subject *domain.Subject) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO subjects (id, user_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		subject.ID, subject.UserID, subject.Name, subject.CreatedAt, subject.UpdatedAt,
	)
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
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, name, created_at, updated_at FROM subjects WHERE id = $1`,
		id,
	)
	var s domain.Subject
	err := row.Scan(&s.ID, &s.UserID, &s.Name, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("subject")
		}
		return nil, fmt.Errorf("find subject: %w", err)
	}
	return domain.ReconstructSubject(s.ID, s.UserID, s.Name, s.CreatedAt, s.UpdatedAt), nil
}

func (r *subjectRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Subject, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, created_at, updated_at FROM subjects WHERE user_id = $1 ORDER BY name`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("find subjects: %w", err)
	}
	defer rows.Close()

	var subjects []*domain.Subject
	for rows.Next() {
		var s domain.Subject
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan subject: %w", err)
		}
		subjects = append(subjects, domain.ReconstructSubject(s.ID, s.UserID, s.Name, s.CreatedAt, s.UpdatedAt))
	}
	return subjects, rows.Err()
}

func (r *subjectRepository) Update(ctx context.Context, subject *domain.Subject) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE subjects SET name = $1, updated_at = $2 WHERE id = $3`,
		subject.Name, subject.UpdatedAt, subject.ID,
	)
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
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM subjects WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete subject: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("subject")
	}
	return nil
}
