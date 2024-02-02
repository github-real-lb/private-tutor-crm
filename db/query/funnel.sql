-- name: GetFunnel :one
SELECT * FROM funnels
WHERE funnel_id = $1 LIMIT 1;

-- name: ListFunnels :many
SELECT * FROM funnels
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: CreateFunnel :one
INSERT INTO funnels (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdateFunnel :exec
UPDATE funnels
  set name = $2
WHERE funnel_id = $1;

-- name: DeleteFunnel :exec
DELETE FROM funnels
WHERE funnel_id = $1;