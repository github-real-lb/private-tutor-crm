package db

import (
	"context"
	"database/sql"
	"sort"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomPayment adds a new random payment, and returns the Payment data type.
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

	require.NotZero(t, payment.PaymentID)

	require.Equal(t, arg.ReceiptID, payment.ReceiptID)
	require.WithinDuration(t, arg.PaymentDatetime, payment.PaymentDatetime, time.Second)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.PaymentMethodID, payment.PaymentMethodID)

	return payment
}

// createRandomPayments adds 'n' random payments with the same ReceiptID, and returns the Payments type.
func createRandomPayments(t *testing.T, n int) Payments {
	var payments Payments

	receipt := createRandomReceipt(t)
	paymentMethod := createRandomPaymentMethod(t)

	for i := 0; i < n; i++ {
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

		payments = append(payments, payment)
	}

	return payments
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

func TestGetPayments(t *testing.T) {
	n := 5 // number of payments to create
	payments1 := createRandomPayments(t, n)

	payments2, err := testQueries.GetPayments(context.Background(), payments1[0].ReceiptID)
	require.NoError(t, err)
	require.NotEmpty(t, payments2)
	require.Equal(t, len(payments2), n)

	sort.Sort(payments1)

	for i := 0; i < n; i++ {
		payment1 := payments1[i]
		payment2 := payments2[i]

		require.NotEmpty(t, payment2)

		require.Equal(t, payment1.PaymentID, payment2.PaymentID)
		require.Equal(t, payment1.ReceiptID, payment2.ReceiptID)
		require.WithinDuration(t, payment1.PaymentDatetime, payment2.PaymentDatetime, time.Second)
		require.Equal(t, payment1.Amount, payment2.Amount)
		require.Equal(t, payment1.PaymentMethodID, payment2.PaymentMethodID)
	}
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
	n := 5 // number of payments to create

	payments := createRandomPayments(t, n)

	err := testQueries.DeletePaymentsByReceipt(context.Background(), payments[0].ReceiptID)
	require.NoError(t, err)

	for _, v := range payments {
		payment, err := testQueries.GetPayment(context.Background(), v.PaymentID)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, payment)
	}
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
