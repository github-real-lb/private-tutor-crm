-- name: GetLessonLocation :one
SELECT * FROM lesson_locations
WHERE location_id = $1 LIMIT 1;

-- name: ListLessonLocations :many
SELECT * FROM lesson_locations
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: CreateLessonLocation :one
INSERT INTO lesson_locations (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdateLessonLocation :exec
UPDATE lesson_locations
  set name = $2
WHERE location_id = $1;

-- name: DeleteLessonLocation :exec
DELETE FROM lesson_locations
WHERE location_id = $1;