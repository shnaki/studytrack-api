package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/repository/sqlcgen"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

type studyLogRepository struct {
	q    *sqlcgen.Queries
	pool *pgxpool.Pool
}

func NewStudyLogRepository(pool *pgxpool.Pool) port.StudyLogRepository {
	return &studyLogRepository{
		q:    sqlcgen.New(pool),
		pool: pool,
	}
}

func (r *studyLogRepository) Create(ctx context.Context, log *domain.StudyLog) error {
	err := r.q.CreateStudyLog(ctx, sqlcgen.CreateStudyLogParams{
		ID:        toPgUUID(log.ID),
		UserID:    toPgUUID(log.UserID),
		SubjectID: toPgUUID(log.SubjectID),
		StudiedAt: toPgTimestamptz(log.StudiedAt),
		Minutes:   int32(log.Minutes),
		Note:      log.Note,
		CreatedAt: toPgTimestamptz(log.CreatedAt),
	})
	if err != nil {
		return fmt.Errorf("insert study log: %w", err)
	}
	return nil
}

func (r *studyLogRepository) FindByID(ctx context.Context, id string) (*domain.StudyLog, error) {
	row, err := r.q.GetStudyLogByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("study log")
		}
		return nil, fmt.Errorf("find study log: %w", err)
	}
	return domain.ReconstructStudyLog(
		fromPgUUID(row.ID),
		fromPgUUID(row.UserID),
		fromPgUUID(row.SubjectID),
		fromPgTimestamptz(row.StudiedAt),
		int(row.Minutes),
		row.Note,
		fromPgTimestamptz(row.CreatedAt),
	), nil
}

// FindByUserID uses dynamic SQL for flexible filtering, so it bypasses sqlcgen.
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
	tag, err := r.q.DeleteStudyLog(ctx, toPgUUID(id))
	if err != nil {
		return fmt.Errorf("delete study log: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound("study log")
	}
	return nil
}
