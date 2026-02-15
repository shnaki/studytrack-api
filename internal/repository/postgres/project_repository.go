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

type projectRepository struct {
	q *sqlcgen.Queries
}

// NewProjectRepository creates a new ProjectRepository implementation using PostgreSQL.
func NewProjectRepository(pool *pgxpool.Pool) port.ProjectRepository {
	return &projectRepository{q: sqlcgen.New(pool)}
}

func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	err := r.q.CreateProject(ctx, sqlcgen.CreateProjectParams{
		ID:        toPgUUID(project.ID),
		UserID:    toPgUUID(project.UserID),
		Name:      project.Name,
		CreatedAt: toPgTimestamptz(project.CreatedAt),
		UpdatedAt: toPgTimestamptz(project.UpdatedAt),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrConflict("project with this name already exists for this user")
		}
		return fmt.Errorf("insert project: %w", err)
	}
	return nil
}

func (r *projectRepository) FindByID(ctx context.Context, id string) (*domain.Project, error) {
	row, err := r.q.GetProjectByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("project")
		}
		return nil, fmt.Errorf("find project: %w", err)
	}
	return domain.ReconstructProject(
		fromPgUUID(row.ID),
		fromPgUUID(row.UserID),
		row.Name,
		fromPgTimestamptz(row.CreatedAt),
		fromPgTimestamptz(row.UpdatedAt),
	), nil
}

func (r *projectRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Project, error) {
	rows, err := r.q.ListProjectsByUserID(ctx, toPgUUID(userID))
	if err != nil {
		return nil, fmt.Errorf("find projects: %w", err)
	}
	projects := make([]*domain.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, domain.ReconstructProject(
			fromPgUUID(row.ID),
			fromPgUUID(row.UserID),
			row.Name,
			fromPgTimestamptz(row.CreatedAt),
			fromPgTimestamptz(row.UpdatedAt),
		))
	}
	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	tag, err := r.q.UpdateProject(ctx, sqlcgen.UpdateProjectParams{
		Name:      project.Name,
		UpdatedAt: toPgTimestamptz(project.UpdatedAt),
		ID:        toPgUUID(project.ID),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrConflict("project with this name already exists for this user")
		}
		return fmt.Errorf("update project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("project")
	}
	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.q.DeleteProject(ctx, toPgUUID(id))
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("project")
	}
	return nil
}
