-- name: GetPaymentMethod :one
SELECT * FROM payment_methods
WHERE payment_method_id = $1 LIMIT 1;

-- name: ListPaymentMethods :many
SELECT * FROM payment_methods
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: CreatePaymentMethod :one
INSERT INTO payment_methods (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdatePaymentMethod :exec
UPDATE payment_methods
  set name = $2
WHERE payment_method_id = $1;

-- name: DeletePaymentMethod :exec
DELETE FROM payment_methods
WHERE payment_method_id = $1;