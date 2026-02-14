package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

type studyLogRepository struct {
	pool *pgxpool.Pool
}

func NewStudyLogRepository(pool *pgxpool.Pool) port.StudyLogRepository {
	return &studyLogRepository{pool: pool}
}

func (r *studyLogRepository) Create(ctx context.Context, log *domain.StudyLog) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO study_logs (id, user_id, subject_id, studied_at, minutes, note, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		log.ID, log.UserID, log.SubjectID, log.StudiedAt, log.Minutes, log.Note, log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert study log: %w", err)
	}
	return nil
}

func (r *studyLogRepository) FindByID(ctx context.Context, id string) (*domain.StudyLog, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, subject_id, studied_at, minutes, note, created_at
		 FROM study_logs WHERE id = $1`,
		id,
	)
	var l domain.StudyLog
	err := row.Scan(&l.ID, &l.UserID, &l.SubjectID, &l.StudiedAt, &l.Minutes, &l.Note, &l.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("study log")
		}
		return nil, fmt.Errorf("find study log: %w", err)
	}
	return domain.ReconstructStudyLog(l.ID, l.UserID, l.SubjectID, l.StudiedAt, l.Minutes, l.Note, l.CreatedAt), nil
}

func (r *studyLogRepository) FindByUserID(ctx context.Context, userID string, filter port.StudyLogFilter) ([]*domain.StudyLog, error) {
	query := strings.Builder{}
	query.WriteString(`SELECT id, user_id, subject_id, studied_at, minutes, note, created_at FROM study_logs WHERE user_id = $1`)

	args := []any{userID}
	paramIdx := 2

	if filter.From != nil {
		query.WriteString(fmt.Sprintf(` AND studied_at >= $%d`, paramIdx))
		args = append(args, *filter.From)
		paramIdx++
	}
	if filter.To != nil {
		query.WriteString(fmt.Sprintf(` AND studied_at < $%d`, paramIdx))
		args = append(args, *filter.To)
		paramIdx++
	}
	if filter.SubjectID != nil {
		query.WriteString(fmt.Sprintf(` AND subject_id = $%d`, paramIdx))
		args = append(args, *filter.SubjectID)
	}

	query.WriteString(` ORDER BY studied_at DESC`)

	rows, err := r.pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("find study logs: %w", err)
	}
	defer rows.Close()

	var logs []*domain.StudyLog
	for rows.Next() {
		var l domain.StudyLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.SubjectID, &l.StudiedAt, &l.Minutes, &l.Note, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan study log: %w", err)
		}
		logs = append(logs, domain.ReconstructStudyLog(l.ID, l.UserID, l.SubjectID, l.StudiedAt, l.Minutes, l.Note, l.CreatedAt))
	}
	return logs, rows.Err()
}

func (r *studyLogRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM study_logs WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete study log: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("study log")
	}
	return nil
}
