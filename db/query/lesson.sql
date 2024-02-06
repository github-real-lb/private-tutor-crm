-- name: CreateLesson :one
INSERT INTO lessons (
  lesson_datetime, duration, location_id, subject_id, notes
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetLesson :one
SELECT * FROM lessons
WHERE lesson_id = $1 LIMIT 1;

-- name: ListLessons :many
SELECT * FROM lessons
ORDER BY lesson_datetime
LIMIT $1
OFFSET $2;

-- name: UpdateLesson :exec
UPDATE lessons
  set   lesson_datetime = $2, 
        duration = $3,
        location_id = $4, 
        subject_id =  $5,
        notes = $6
WHERE lesson_id = $1;

-- name: DeleteLesson :exec
DELETE FROM lessons
WHERE lesson_id = $1;