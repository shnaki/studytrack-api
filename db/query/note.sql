-- name: CreateNote :exec
INSERT INTO notes (id, project_id, user_id, title, content, tags, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetNoteByID :one
SELECT id, project_id, user_id, title, content, tags, created_at, updated_at
FROM notes
WHERE id = $1;

-- name: ListNotesByProjectID :many
SELECT id, project_id, user_id, title, content, tags, created_at, updated_at
FROM notes
WHERE project_id = $1
ORDER BY updated_at DESC;

-- name: UpdateNote :execresult
UPDATE notes SET title = $1, content = $2, tags = $3, updated_at = $4 WHERE id = $5;

-- name: DeleteNote :execresult
DELETE FROM notes WHERE id = $1;
