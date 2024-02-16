package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	CreateReceiptWithPaymentsTx(ctx context.Context, arg CreateReceiptTxParams) (ReceiptWithPayments, error)
	GetReceiptWithPaymentsTx(ctx context.Context, receiptID int64) (ReceiptWithPayments, error)
	DeleteReceiptWithPaymentsTx(ctx context.Context, receiptID int64) error
	GetReceiptsWithPaymentsByStudentTx(ctx context.Context, studentID int64, limit, offset int) (StudentReceiptsWithPayments, error)
	CreateLessonWithInvoicesTx(ctx context.Context, arg CreateLessonTxParams) (LessonWithInvoices, error)
	GetLessonWithInvoicesTx(ctx context.Context, lessonID int64) (LessonWithInvoices, error)
	DeleteLessonWithInvoicesTx(ctx context.Context, lessonID int64) error
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
