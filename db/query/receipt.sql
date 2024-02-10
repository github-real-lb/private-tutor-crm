-- name: CreateReceipt :one
INSERT INTO receipts (
  student_id, receipt_datetime, amount, notes
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetReceipt :one
SELECT * FROM receipts
WHERE receipt_id = $1 LIMIT 1;

-- name: ListReceipts :many
SELECT * FROM receipts
ORDER BY student_id, receipt_datetime
LIMIT $1
OFFSET $2;

-- name: UpdateReceipt :exec
UPDATE receipts
  set   student_id = $2,
        receipt_datetime = $3, 
        amount = $4,
        notes = $5
WHERE receipt_id = $1;

-- name: UpdateReceiptAmount :exec
UPDATE receipts
  set   amount = $2
WHERE receipt_id = $1;

-- name: DeleteReceipt :exec
DELETE FROM receipts
WHERE receipt_id = $1;