package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/repository/sqlcgen"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

type noteRepository struct {
	q *sqlcgen.Queries
}

// NewNoteRepository creates a new NoteRepository implementation using PostgreSQL.
func NewNoteRepository(pool *pgxpool.Pool) port.NoteRepository {
	return &noteRepository{q: sqlcgen.New(pool)}
}

func (r *noteRepository) Create(ctx context.Context, note *domain.Note) error {
	err := r.q.CreateNote(ctx, sqlcgen.CreateNoteParams{
		ID:        toPgUUID(note.ID),
		ProjectID: toPgUUID(note.ProjectID),
		UserID:    toPgUUID(note.UserID),
		Title:     note.Title,
		Content:   note.Content,
		Tags:      note.Tags,
		CreatedAt: toPgTimestamptz(note.CreatedAt),
		UpdatedAt: toPgTimestamptz(note.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("insert note: %w", err)
	}
	return nil
}

func (r *noteRepository) FindByID(ctx context.Context, id string) (*domain.Note, error) {
	row, err := r.q.GetNoteByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("note")
		}
		return nil, fmt.Errorf("find note: %w", err)
	}
	return domain.ReconstructNote(
		fromPgUUID(row.ID),
		fromPgUUID(row.ProjectID),
		fromPgUUID(row.UserID),
		row.Title,
		row.Content,
		row.Tags,
		fromPgTimestamptz(row.CreatedAt),
		fromPgTimestamptz(row.UpdatedAt),
	), nil
}

func (r *noteRepository) FindByProjectID(ctx context.Context, projectID string) ([]*domain.Note, error) {
	rows, err := r.q.ListNotesByProjectID(ctx, toPgUUID(projectID))
	if err != nil {
		return nil, fmt.Errorf("find notes: %w", err)
	}
	notes := make([]*domain.Note, 0, len(rows))
	for _, row := range rows {
		notes = append(notes, domain.ReconstructNote(
			fromPgUUID(row.ID),
			fromPgUUID(row.ProjectID),
			fromPgUUID(row.UserID),
			row.Title,
			row.Content,
			row.Tags,
			fromPgTimestamptz(row.CreatedAt),
			fromPgTimestamptz(row.UpdatedAt),
		))
	}
	return notes, nil
}

func (r *noteRepository) Update(ctx context.Context, note *domain.Note) error {
	tag, err := r.q.UpdateNote(ctx, sqlcgen.UpdateNoteParams{
		Title:     note.Title,
		Content:   note.Content,
		Tags:      note.Tags,
		UpdatedAt: toPgTimestamptz(note.UpdatedAt),
		ID:        toPgUUID(note.ID),
	})
	if err != nil {
		return fmt.Errorf("update note: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("note")
	}
	return nil
}

func (r *noteRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.q.DeleteNote(ctx, toPgUUID(id))
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("note")
	}
	return nil
}
