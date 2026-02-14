package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/application/port"
	"github.com/shnaki/studytrack-api/internal/domain"
)

type goalRepository struct {
	pool *pgxpool.Pool
}

func NewGoalRepository(pool *pgxpool.Pool) port.GoalRepository {
	return &goalRepository{pool: pool}
}

func (r *goalRepository) Upsert(ctx context.Context, goal *domain.Goal) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO goals (id, user_id, subject_id, target_minutes_per_week, start_date, end_date, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (user_id, subject_id)
		 DO UPDATE SET target_minutes_per_week = EXCLUDED.target_minutes_per_week,
		               start_date = EXCLUDED.start_date,
		               end_date = EXCLUDED.end_date,
		               updated_at = EXCLUDED.updated_at`,
		goal.ID, goal.UserID, goal.SubjectID, goal.TargetMinutesPerWeek,
		goal.StartDate, goal.EndDate, goal.CreatedAt, goal.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert goal: %w", err)
	}
	return nil
}

func (r *goalRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Goal, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, subject_id, target_minutes_per_week, start_date, end_date, created_at, updated_at
		 FROM goals WHERE user_id = $1 ORDER BY created_at`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("find goals: %w", err)
	}
	defer rows.Close()

	var goals []*domain.Goal
	for rows.Next() {
		var g domain.Goal
		if err := rows.Scan(&g.ID, &g.UserID, &g.SubjectID, &g.TargetMinutesPerWeek,
			&g.StartDate, &g.EndDate, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan goal: %w", err)
		}
		goals = append(goals, domain.ReconstructGoal(
			g.ID, g.UserID, g.SubjectID, g.TargetMinutesPerWeek,
			g.StartDate, g.EndDate, g.CreatedAt, g.UpdatedAt,
		))
	}
	return goals, rows.Err()
}
