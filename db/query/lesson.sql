-- name: CreateLesson :one
INSERT INTO lessons (
  student_id, lesson_datetime, duration, location_id, subject_id, notes
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetLesson :one
SELECT * FROM lessons
WHERE lesson_id = $1 LIMIT 1;

-- name: ListLessons :many
SELECT * FROM lessons
ORDER BY student_id, lesson_datetime
LIMIT $1
OFFSET $2;

-- name: UpdateLesson :exec
UPDATE lessons
  set   student_id = $2,
        lesson_datetime = $3, 
        duration = $4,
        location_id = $5, 
        subject_id =  $6,
        notes = $7
WHERE lesson_id = $1;

-- name: DeleteLesson :exec
DELETE FROM lessons
WHERE lesson_id = $1;