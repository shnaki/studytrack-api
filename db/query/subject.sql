-- name: CreateSubject :exec
INSERT INTO subjects (id, user_id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetSubjectByID :one
SELECT id, user_id, name, created_at, updated_at
FROM subjects
WHERE id = $1;

-- name: ListSubjectsByUserID :many
SELECT id, user_id, name, created_at, updated_at
FROM subjects
WHERE user_id = $1
ORDER BY name;

-- name: UpdateSubject :exec
UPDATE subjects SET name = $1, updated_at = $2 WHERE id = $3;

-- name: DeleteSubject :exec
DELETE FROM subjects WHERE id = $1;
