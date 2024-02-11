-- name: CreateInvoice :one
INSERT INTO invoices (
  student_id, lesson_id, invoice_datetime, hourly_fee, duration, discount, amount, notes
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices
WHERE invoice_id = $1 LIMIT 1;

-- name: GetInvoicesByLesson :many
SELECT * FROM invoices
WHERE lesson_id = $1
ORDER BY student_id;

-- name: GetInvoicesByStudent :many
SELECT * FROM invoices
WHERE student_id = $1
ORDER BY invoice_datetime;

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
        duration = $6, 
        discount = $7,
        amount =  $8,
        notes = $9
WHERE invoice_id = $1;

-- name: DeleteInvoice :exec
DELETE FROM invoices
WHERE invoice_id = $1;

-- name: DeleteInvoicesByLesson :exec
DELETE FROM invoices
WHERE lesson_id = $1;