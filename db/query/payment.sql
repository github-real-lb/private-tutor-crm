-- name: CreatePayment :one
INSERT INTO payments (
  receipt_id, payment_datetime, amount, payment_method_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE payment_id = $1 LIMIT 1;

-- name: GetPayments :many
SELECT * FROM payments
WHERE receipt_id = $1 
ORDER BY receipt_id, payment_datetime;

-- name: ListPayments :many
SELECT * FROM payments
ORDER BY receipt_id, payment_datetime
LIMIT $1
OFFSET $2;

-- name: UpdatePayment :exec
UPDATE payments
  set   receipt_id = $2,
        payment_datetime = $3, 
        amount = $4,
        payment_method_id = $5
WHERE payment_id = $1;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE payment_id = $1;

-- name: DeletePayments :exec
DELETE FROM payments
WHERE receipt_id = $1;