-- name: CreateProject :exec
INSERT INTO projects (id, user_id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetProjectByID :one
SELECT id, user_id, name, created_at, updated_at
FROM projects
WHERE id = $1;

-- name: ListProjectsByUserID :many
SELECT id, user_id, name, created_at, updated_at
FROM projects
WHERE user_id = $1
ORDER BY name;

-- name: UpdateProject :execresult
UPDATE projects SET name = $1, updated_at = $2 WHERE id = $3;

-- name: DeleteProject :execresult
DELETE FROM projects WHERE id = $1;
