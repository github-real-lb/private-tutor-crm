package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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

// LessonInvoicingTxInvoiceParams contains the input paramaters of a single invoice in a lesson's invoicing transaction.
type LessonInvoicingTxInvoiceParams struct {
	StudentID int64          `json:"student_id"`
	HourlyFee float64        `json:"hourly_fee"`
	Duration  int64          `json:"duration"`
	Discount  float64        `json:"discount"`
	Amount    float64        `json:"amount"`
	Notes     sql.NullString `json:"notes"`
}

// LessonInvoicingTxParams contains the input paramaters of the lesson's invoicing transaction.
type LessonInvoicingTxParams struct {
	LessonDatetime       time.Time                        `json:"lesson_datetime"`
	Duration             int64                            `json:"duration"`
	LocationID           int64                            `json:"location_id"`
	SubjectID            int64                            `json:"subject_id"`
	Notes                sql.NullString                   `json:"notes"`
	LessonInvoicesParams []LessonInvoicingTxInvoiceParams `json:"lesson_invoices_params"`
}

// LessonInvoicingTxResult is the result of the invoicing transaction.
type LessonInvoicingTxResult struct {
	Lesson   Lesson    `json:"lesson"`
	Invoices []Invoice `json:"invoices"`
}

// LessonInvoicingTx creates a lesson that took place,
// and invoices for all the students that took part in the lesson.
func (store *Store) LessonInvoicingTx(ctx context.Context, arg LessonInvoicingTxParams) (LessonInvoicingTxResult, error) {
	var result LessonInvoicingTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		createLessonArg := CreateLessonParams{
			LessonDatetime: arg.LessonDatetime,
			Duration:       arg.Duration,
			LocationID:     arg.LocationID,
			SubjectID:      arg.SubjectID,
			Notes:          arg.Notes,
		}

		result.Lesson, err = q.CreateLesson(ctx, createLessonArg)
		if err != nil {
			return err
		}

		for _, invoiceArg := range arg.LessonInvoicesParams {
			createInvoiceArg := CreateInvoiceParams{
				StudentID: invoiceArg.StudentID,
				LessonID:  result.Lesson.LessonID,
				HourlyFee: invoiceArg.HourlyFee,
				Duration:  invoiceArg.Duration,
				Discount:  invoiceArg.Discount,
				Amount:    invoiceArg.Amount,
				Notes:     invoiceArg.Notes,
			}

			invoice, err := q.CreateInvoice(ctx, createInvoiceArg)
			if err != nil {
				return err
			}

			result.Invoices = append(result.Invoices, invoice)
		}

		return nil
	})

	return result, err
}

type PaymentsReceivingTxPaymentParams struct {
	PaymentDatetime time.Time `json:"payment_datetime"`
	Amount          float64   `json:"amount"`
	PaymentMethodID int64     `json:"payment_method_id"`
}

type PaymentsReceivingTxParams struct {
	StudentID             int64                              `json:"student_id"`
	ReceiptDatetime       time.Time                          `json:"receipt_datetime"`
	Amount                float64                            `json:"amount"`
	Notes                 sql.NullString                     `json:"notes"`
	ReceiptPaymentsParams []PaymentsReceivingTxPaymentParams `json:"receipt_payments_params"`
}

type PaymentsReceivingTxResult struct {
	Receipt  Receipt   `json:"receipt"`
	Payments []Payment `json:"payments"`
}

func (store *Store) PaymentsReceivingTx(ctx context.Context, arg PaymentsReceivingTxParams) (PaymentsReceivingTxResult, error) {
	var result PaymentsReceivingTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		createReceiptArg := CreateReceiptParams{
			StudentID:       arg.StudentID,
			ReceiptDatetime: arg.ReceiptDatetime,
			Amount:          arg.Amount,
			Notes:           arg.Notes,
		}

		result.Receipt, err = q.CreateReceipt(ctx, createReceiptArg)
		if err != nil {
			return err
		}

		for _, paymentArg := range arg.ReceiptPaymentsParams {
			createPaymentArg := CreatePaymentParams{
				ReceiptID:       result.Receipt.ReceiptID,
				PaymentDatetime: paymentArg.PaymentDatetime,
				Amount:          paymentArg.Amount,
				PaymentMethodID: paymentArg.PaymentMethodID,
			}

			payment, err := q.CreatePayment(ctx, createPaymentArg)
			if err != nil {
				return err
			}

			result.Payments = append(result.Payments, payment)
		}

		return nil
	})

	return result, err
}
