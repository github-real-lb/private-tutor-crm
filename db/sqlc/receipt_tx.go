package db

import (
	"context"
	"database/sql"
	"time"
)

// CreateReceiptTxPaymentParams contains the input paramaters of a single payment, for the CreateReceiptTx function.
type CreateReceiptTxPaymentParams struct {
	PaymentDatetime time.Time `json:"payment_datetime"`
	Amount          float64   `json:"amount"`
	PaymentMethodID int64     `json:"payment_method_id"`
}

// CreateReceiptTxParams contains the input paramaters of a single reciept and its payments, for the CreateReceiptTx function.
type CreateReceiptTxParams struct {
	StudentID             int64                          `json:"student_id"`
	ReceiptDatetime       time.Time                      `json:"receipt_datetime"`
	Notes                 sql.NullString                 `json:"notes"`
	ReceiptPaymentsParams []CreateReceiptTxPaymentParams `json:"receipt_payments_params"`
}

// ReceiptTxResult is the result receipt and its payments, for any lesson transaction.
type ReceiptTxResult struct {
	Receipt  Receipt   `json:"receipt"`
	Payments []Payment `json:"payments"`
}

// CreateReceiptTx creates a Receipt and all the Payments releated to it.
// Receipt amount is calculated end updated based on all payments.
func (store *Store) CreateReceiptTx(ctx context.Context, arg CreateReceiptTxParams) (ReceiptTxResult, error) {
	var result ReceiptTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		createReceiptArg := CreateReceiptParams{
			StudentID:       arg.StudentID,
			ReceiptDatetime: arg.ReceiptDatetime,
			Amount:          0.0,
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
				PaymentMethodID: paymentArg.PaymentMethodID,
			}

			payment, err := q.CreatePayment(ctx, createPaymentArg)
			if err != nil {
				return err
			}

			result.Receipt.Amount += payment.Amount
			result.Payments = append(result.Payments, payment)
		}

		err = q.UpdateReceiptAmount(ctx, UpdateReceiptAmountParams{
			ReceiptID: result.Receipt.ReceiptID,
			Amount:    result.Receipt.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// DeleteReceiptTx deletes a Receipt and all the Payments releated to it.
func (store *Store) DeleteReceiptTx(ctx context.Context, receiptID int64) error {
	err := store.execTx(ctx, func(q *Queries) error {
		err := q.DeletePayments(ctx, receiptID)
		if err != nil {
			return err
		}

		err = q.DeleteReceipt(ctx, receiptID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
