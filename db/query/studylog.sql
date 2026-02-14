-- name: CreateStudyLog :exec
INSERT INTO study_logs (id, user_id, subject_id, studied_at, minutes, note, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetStudyLogByID :one
SELECT id, user_id, subject_id, studied_at, minutes, note, created_at
FROM study_logs
WHERE id = $1;

-- name: DeleteStudyLog :execresult
DELETE FROM study_logs WHERE id = $1;
