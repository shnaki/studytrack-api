package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/repository/sqlcgen"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

type goalRepository struct {
	q *sqlcgen.Queries
}

// NewGoalRepository creates a new GoalRepository implementation using PostgreSQL.
func NewGoalRepository(pool *pgxpool.Pool) port.GoalRepository {
	return &goalRepository{q: sqlcgen.New(pool)}
}

func (r *goalRepository) Upsert(ctx context.Context, goal *domain.Goal) error {
	err := r.q.UpsertGoal(ctx, sqlcgen.UpsertGoalParams{
		ID:                   toPgUUID(goal.ID),
		UserID:               toPgUUID(goal.UserID),
		SubjectID:            toPgUUID(goal.SubjectID),
		TargetMinutesPerWeek: int32(goal.TargetMinutesPerWeek),
		StartDate:            toPgDate(goal.StartDate),
		EndDate:              toPgDatePtr(goal.EndDate),
		CreatedAt:            toPgTimestamptz(goal.CreatedAt),
		UpdatedAt:            toPgTimestamptz(goal.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("upsert goal: %w", err)
	}
	return nil
}

func (r *goalRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Goal, error) {
	rows, err := r.q.ListGoalsByUserID(ctx, toPgUUID(userID))
	if err != nil {
		return nil, fmt.Errorf("find goals: %w", err)
	}
	goals := make([]*domain.Goal, 0, len(rows))
	for _, row := range rows {
		goals = append(goals, domain.ReconstructGoal(
			fromPgUUID(row.ID),
			fromPgUUID(row.UserID),
			fromPgUUID(row.SubjectID),
			int(row.TargetMinutesPerWeek),
			row.StartDate.Time,
			fromPgDatePtr(row.EndDate),
			fromPgTimestamptz(row.CreatedAt),
			fromPgTimestamptz(row.UpdatedAt),
		))
	}
	return goals, nil
}
