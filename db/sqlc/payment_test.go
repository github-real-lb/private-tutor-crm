package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomPayment tests adding a new random payment to the database, and returns the Payment data type.
func createRandomPayment(t *testing.T) Payment {
	receipt := createRandomReceipt(t)
	paymentMethod := createRandomPaymentMethod(t)

	arg := CreatePaymentParams{
		ReceiptID:       receipt.ReceiptID,
		PaymentDatetime: util.RandomDatetime(),
		Amount:          util.RandomPaymentAmount(),
		PaymentMethodID: paymentMethod.PaymentMethodID,
	}

	payment, err := testQueries.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, arg.ReceiptID, payment.ReceiptID)
	require.WithinDuration(t, arg.PaymentDatetime, payment.PaymentDatetime, time.Second)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.PaymentMethodID, payment.PaymentMethodID)

	require.NotZero(t, payment.PaymentID)

	return payment
}

func TestCreatePayment(t *testing.T) {
	createRandomPayment(t)
}

func TestGetPayment(t *testing.T) {
	payment1 := createRandomPayment(t)
	payment2, err := testQueries.GetPayment(context.Background(), payment1.PaymentID)
	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, payment1.PaymentID, payment2.PaymentID)
	require.Equal(t, payment1.ReceiptID, payment2.ReceiptID)
	require.WithinDuration(t, payment1.PaymentDatetime, payment2.PaymentDatetime, time.Second)
	require.Equal(t, payment1.Amount, payment2.Amount)
	require.Equal(t, payment1.PaymentMethodID, payment2.PaymentMethodID)
}

func TestUpdatePayment(t *testing.T) {
	payment1 := createRandomPayment(t)
	receipt := createRandomReceipt(t)
	paymentMethod := createRandomPaymentMethod(t)

	arg := UpdatePaymentParams{
		PaymentID:       payment1.PaymentID,
		ReceiptID:       receipt.ReceiptID,
		PaymentDatetime: util.RandomDatetime(),
		Amount:          util.RandomPaymentAmount(),
		PaymentMethodID: paymentMethod.PaymentMethodID,
	}
	err := testQueries.UpdatePayment(context.Background(), arg)
	require.NoError(t, err)

	payment2, err := testQueries.GetPayment(context.Background(), arg.PaymentID)
	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, arg.PaymentID, payment2.PaymentID)
	require.Equal(t, arg.ReceiptID, payment2.ReceiptID)
	require.WithinDuration(t, arg.PaymentDatetime, payment2.PaymentDatetime, time.Second)
	require.Equal(t, arg.Amount, payment2.Amount)
	require.Equal(t, arg.PaymentMethodID, payment2.PaymentMethodID)
}

func TestDeletePayment(t *testing.T) {
	payment1 := createRandomPayment(t)
	err := testQueries.DeletePayment(context.Background(), payment1.PaymentID)
	require.NoError(t, err)

	payment2, err := testQueries.GetPayment(context.Background(), payment1.PaymentID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, payment2)
}

func TestDeletePayments(t *testing.T) {
	payment1 := createRandomPayment(t)
	payment2 := createRandomPayment(t)

	// update payment2 to have the same receiptID as payment1
	arg := UpdatePaymentParams{
		PaymentID:       payment2.PaymentID,
		ReceiptID:       payment1.ReceiptID,
		PaymentDatetime: payment2.PaymentDatetime,
		Amount:          payment2.Amount,
		PaymentMethodID: payment2.PaymentMethodID,
	}
	err := testQueries.UpdatePayment(context.Background(), arg)
	require.NoError(t, err)

	testQueries.DeletePayments(context.Background(), payment1.ReceiptID)

	payment, err := testQueries.GetPayment(context.Background(), payment1.PaymentID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, payment)

	payment, err = testQueries.GetPayment(context.Background(), payment1.PaymentID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, payment)
}

func TestListPayments(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPayment(t)
	}

	arg := ListPaymentsParams{
		Limit:  5,
		Offset: 5,
	}

	payments, err := testQueries.ListPayments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, payments, 5)

	for _, payment := range payments {
		require.NotEmpty(t, payment)
	}
}
