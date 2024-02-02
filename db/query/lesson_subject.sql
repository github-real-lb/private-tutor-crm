-- name: GetLessonSubject :one
SELECT * FROM lesson_subjects
WHERE subject_id = $1 LIMIT 1;

-- name: ListLessonSubjects :many
SELECT * FROM lesson_subjects
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: CreateLessonSubject :one
INSERT INTO lesson_subjects (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdateLessonSubject :exec
UPDATE lesson_subjects
  set name = $2
WHERE subject_id = $1;

-- name: DeleteLessonSubject :exec
DELETE FROM lesson_subjects
WHERE subject_id = $1;