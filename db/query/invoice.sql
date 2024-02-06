-- name: CreateInvoice :one
INSERT INTO invoices (
  student_id, lesson_id, hourly_fee, duration, discount, amount, notes
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices
WHERE invoice_id = $1 LIMIT 1;

-- name: ListInvoices :many
SELECT * FROM invoices
ORDER BY student_id, created_at
LIMIT $1
OFFSET $2;

-- name: UpdateInvoice :exec
UPDATE invoices
  set   student_id = $2,
        lesson_id = $3, 
        hourly_fee = $4,
        duration = $5, 
        discount = $6,
        amount =  $7,
        notes = $8
WHERE invoice_id = $1;

-- name: DeleteInvoice :exec
DELETE FROM invoices
WHERE invoice_id = $1;