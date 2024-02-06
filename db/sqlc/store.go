package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// InvoiceTxParams contains the input paramaters of the invoicing transaction.
type InvoicingTxParams struct {
	LessonParams   CreateLessonParams    `json:"lesson_params"`
	InvoicesParams []CreateInvoiceParams `json:"invoices_params"`
}

// InvoicingTxResult is the result of the invoicing transaction.
type InvoicingTxResult struct {
	Lesson   Lesson    `json:"lesson"`
	Invoices []Invoice `json:"invoices"`
}

// InvoicingTx creates a lesson that took place,
// and invoices for all the students that participate in the lesson.
func (store *Store) InvoicingTx(ctx context.Context, arg InvoicingTxParams) (InvoicingTxResult, error) {
	var result InvoicingTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Lesson, err = q.CreateLesson(ctx, arg.LessonParams)
		if err != nil {
			return err
		}

		return nil
	})

	for _, invoiceParams := range arg.InvoicesParams {
		err = store.execTx(ctx, func(q *Queries) error {
			var err error

			invoice, err := q.CreateInvoice(ctx, invoiceParams)
			if err != nil {
				return err
			}

			result.Invoices = append(result.Invoices, invoice)

			return nil
		})
	}

	return result, err
}
