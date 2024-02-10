package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomReceiptTx adds a new random receipt with 'n' payments to the database, and returns the ReceiptTxResult data type.
func createRandomReceiptTx(t *testing.T, n int) ReceiptTxResult {
	store := NewStore(testDB)
	student := createRandomStudent(t)
	paymentMethod := createRandomPaymentMethod(t)

	arg := CreateReceiptTxParams{
		StudentID:       student.StudentID,
		ReceiptDatetime: util.RandomDatetime(),
		Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
	}

	for i := 0; i < n; i++ {
		paymentArg := CreateReceiptTxPaymentParams{
			PaymentDatetime: arg.ReceiptDatetime,
			Amount:          util.RandomPaymentAmount(),
			PaymentMethodID: paymentMethod.PaymentMethodID,
		}

		arg.ReceiptPaymentsParams = append(arg.ReceiptPaymentsParams, paymentArg)
	}

	result, err := store.CreateReceiptTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// check Receipt
	require.NotEmpty(t, result.Receipt)
	require.NotZero(t, result.Receipt.ReceiptID)

	receipt, err := testQueries.GetReceipt(context.Background(), result.Receipt.ReceiptID)
	require.NoError(t, err)

	require.Equal(t, receipt.ReceiptID, result.Receipt.ReceiptID)
	require.Equal(t, receipt.StudentID, result.Receipt.StudentID)
	require.WithinDuration(t, receipt.ReceiptDatetime, result.Receipt.ReceiptDatetime, time.Second)
	require.Equal(t, receipt.Notes, result.Receipt.Notes)

	// check Payments
	require.NotEmpty(t, result.Payments)
	require.Equal(t, len(result.Payments), n)

	// check each Payment in Payments
	amount := 0.0
	for _, v := range result.Payments {
		require.NotEmpty(t, v)

		payment, err := testQueries.GetPayment(context.Background(), v.PaymentID)
		require.NoError(t, err)

		require.Equal(t, payment.PaymentID, v.PaymentID)
		require.Equal(t, payment.ReceiptID, v.ReceiptID)
		require.WithinDuration(t, payment.PaymentDatetime, v.PaymentDatetime, time.Second)
		require.Equal(t, payment.Amount, v.Amount)
		require.Equal(t, payment.PaymentMethodID, v.PaymentMethodID)

		amount += payment.Amount
	}

	require.Equal(t, receipt.Amount, amount)

	return result
}

func TestCreateReceiptTx(t *testing.T) {
	createRandomReceiptTx(t, 2)
}

func TestDeleteReceiptTx(t *testing.T) {
	store := NewStore(testDB)

	result := createRandomReceiptTx(t, 2)
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Receipt)
	require.NotZero(t, result.Receipt.ReceiptID)

	err := store.DeleteReceiptTx(context.Background(), result.Receipt.ReceiptID)
	require.NoError(t, err)

	// check receipt deleted
	receipt, err := testQueries.GetReceipt(context.Background(), result.Receipt.ReceiptID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, receipt)

	// check payments deleted
	for _, v := range result.Payments {
		payment, err := testQueries.GetPayment(context.Background(), v.PaymentID)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, payment)
	}

}
