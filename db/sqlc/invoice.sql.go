// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: invoice.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createInvoice = `-- name: CreateInvoice :one
INSERT INTO invoices (
  student_id, lesson_id, invoice_datetime, hourly_fee, duration, discount, amount, notes
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING invoice_id, student_id, lesson_id, invoice_datetime, hourly_fee, duration, discount, amount, notes
`

type CreateInvoiceParams struct {
	StudentID       int64          `json:"student_id"`
	LessonID        int64          `json:"lesson_id"`
	InvoiceDatetime time.Time      `json:"invoice_datetime"`
	HourlyFee       float64        `json:"hourly_fee"`
	Duration        int64          `json:"duration"`
	Discount        float64        `json:"discount"`
	Amount          float64        `json:"amount"`
	Notes           sql.NullString `json:"notes"`
}

func (q *Queries) CreateInvoice(ctx context.Context, arg CreateInvoiceParams) (*Invoice, error) {
	row := q.db.QueryRowContext(ctx, createInvoice,
		arg.StudentID,
		arg.LessonID,
		arg.InvoiceDatetime,
		arg.HourlyFee,
		arg.Duration,
		arg.Discount,
		arg.Amount,
		arg.Notes,
	)
	var i Invoice
	err := row.Scan(
		&i.InvoiceID,
		&i.StudentID,
		&i.LessonID,
		&i.InvoiceDatetime,
		&i.HourlyFee,
		&i.Duration,
		&i.Discount,
		&i.Amount,
		&i.Notes,
	)
	return &i, err
}

const deleteInvoice = `-- name: DeleteInvoice :exec
DELETE FROM invoices
WHERE invoice_id = $1
`

func (q *Queries) DeleteInvoice(ctx context.Context, invoiceID int64) error {
	_, err := q.db.ExecContext(ctx, deleteInvoice, invoiceID)
	return err
}

const deleteInvoicesByLesson = `-- name: DeleteInvoicesByLesson :exec
DELETE FROM invoices
WHERE lesson_id = $1
`

func (q *Queries) DeleteInvoicesByLesson(ctx context.Context, lessonID int64) error {
	_, err := q.db.ExecContext(ctx, deleteInvoicesByLesson, lessonID)
	return err
}

const getInvoice = `-- name: GetInvoice :one
SELECT invoice_id, student_id, lesson_id, invoice_datetime, hourly_fee, duration, discount, amount, notes FROM invoices
WHERE invoice_id = $1 LIMIT 1
`

func (q *Queries) GetInvoice(ctx context.Context, invoiceID int64) (*Invoice, error) {
	row := q.db.QueryRowContext(ctx, getInvoice, invoiceID)
	var i Invoice
	err := row.Scan(
		&i.InvoiceID,
		&i.StudentID,
		&i.LessonID,
		&i.InvoiceDatetime,
		&i.HourlyFee,
		&i.Duration,
		&i.Discount,
		&i.Amount,
		&i.Notes,
	)
	return &i, err
}

const getInvoicesByLesson = `-- name: GetInvoicesByLesson :many
SELECT invoice_id, student_id, lesson_id, invoice_datetime, hourly_fee, duration, discount, amount, notes FROM invoices
WHERE lesson_id = $1
ORDER BY student_id
`

func (q *Queries) GetInvoicesByLesson(ctx context.Context, lessonID int64) ([]*Invoice, error) {
	rows, err := q.db.QueryContext(ctx, getInvoicesByLesson, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.InvoiceID,
			&i.StudentID,
			&i.LessonID,
			&i.InvoiceDatetime,
			&i.HourlyFee,
			&i.Duration,
			&i.Discount,
			&i.Amount,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInvoicesByStudent = `-- name: GetInvoicesByStudent :many
SELECT invoice_id, student_id, lesson_id, invoice_datetime, hourly_fee, duration, discount, amount, notes FROM invoices
WHERE student_id = $1
ORDER BY invoice_datetime
`

func (q *Queries) GetInvoicesByStudent(ctx context.Context, studentID int64) ([]*Invoice, error) {
	rows, err := q.db.QueryContext(ctx, getInvoicesByStudent, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.InvoiceID,
			&i.StudentID,
			&i.LessonID,
			&i.InvoiceDatetime,
			&i.HourlyFee,
			&i.Duration,
			&i.Discount,
			&i.Amount,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listInvoices = `-- name: ListInvoices :many
SELECT invoice_id, student_id, lesson_id, invoice_datetime, hourly_fee, duration, discount, amount, notes FROM invoices
ORDER BY student_id, invoice_datetime
LIMIT $1
OFFSET $2
`

type ListInvoicesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListInvoices(ctx context.Context, arg ListInvoicesParams) ([]*Invoice, error) {
	rows, err := q.db.QueryContext(ctx, listInvoices, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(
			&i.InvoiceID,
			&i.StudentID,
			&i.LessonID,
			&i.InvoiceDatetime,
			&i.HourlyFee,
			&i.Duration,
			&i.Discount,
			&i.Amount,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateInvoice = `-- name: UpdateInvoice :exec
UPDATE invoices
  set   student_id = $2,
        lesson_id = $3, 
        invoice_datetime = $4,
        hourly_fee = $5,
        duration = $6, 
        discount = $7,
        amount =  $8,
        notes = $9
WHERE invoice_id = $1
`

type UpdateInvoiceParams struct {
	InvoiceID       int64          `json:"invoice_id"`
	StudentID       int64          `json:"student_id"`
	LessonID        int64          `json:"lesson_id"`
	InvoiceDatetime time.Time      `json:"invoice_datetime"`
	HourlyFee       float64        `json:"hourly_fee"`
	Duration        int64          `json:"duration"`
	Discount        float64        `json:"discount"`
	Amount          float64        `json:"amount"`
	Notes           sql.NullString `json:"notes"`
}

func (q *Queries) UpdateInvoice(ctx context.Context, arg UpdateInvoiceParams) error {
	_, err := q.db.ExecContext(ctx, updateInvoice,
		arg.InvoiceID,
		arg.StudentID,
		arg.LessonID,
		arg.InvoiceDatetime,
		arg.HourlyFee,
		arg.Duration,
		arg.Discount,
		arg.Amount,
		arg.Notes,
	)
	return err
}
