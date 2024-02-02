-- name: GetCollege :one
SELECT * FROM colleges
WHERE college_id = $1 LIMIT 1;

-- name: ListColleges :many
SELECT * FROM colleges
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: CreateCollege :one
INSERT INTO colleges (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdateCollege :exec
UPDATE colleges
  set name = $2
WHERE college_id = $1;

-- name: DeleteCollege :exec
DELETE FROM colleges
WHERE college_id = $1;