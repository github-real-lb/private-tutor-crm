-- name: CreateInvoice :one
INSERT INTO invoices (
  student_id, lesson_id, invoice_datetime, hourly_fee, amount, notes
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices
WHERE invoice_id = $1 LIMIT 1;

-- name: ListInvoices :many
SELECT * FROM invoices
ORDER BY student_id, invoice_datetime
LIMIT $1
OFFSET $2;

-- name: UpdateInvoice :exec
UPDATE invoices
  set   student_id = $2,
        lesson_id = $3, 
        invoice_datetime = $4,
        hourly_fee = $5, 
        amount =  $6,
        notes = $7
WHERE invoice_id = $1;

-- name: DeleteInvoice :exec
DELETE FROM invoices
WHERE invoice_id = $1;