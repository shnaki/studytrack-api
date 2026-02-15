-- name: UpsertGoal :exec
INSERT INTO goals (id, user_id, project_id, target_minutes_per_week, start_date, end_date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (user_id, project_id)
DO UPDATE SET target_minutes_per_week = EXCLUDED.target_minutes_per_week,
              start_date = EXCLUDED.start_date,
              end_date = EXCLUDED.end_date,
              updated_at = EXCLUDED.updated_at;

-- name: ListGoalsByUserID :many
SELECT id, user_id, project_id, target_minutes_per_week, start_date, end_date, created_at, updated_at
FROM goals
WHERE user_id = $1
ORDER BY created_at;
