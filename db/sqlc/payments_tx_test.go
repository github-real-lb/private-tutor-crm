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

// createRandomReceiptWithPaymentsTx adds a new random receipt with 'n' payments to the database.
func createRandomReceiptWithPaymentsTx(t *testing.T, nPayments int) ReceiptWithPayments {
	store := NewStore(testDB)
	student := createRandomStudent(t)
	paymentMethod := createRandomPaymentMethod(t)

	arg := CreateReceiptTxParams{
		StudentID:       student.StudentID,
		ReceiptDatetime: util.RandomDatetime(),
		Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
	}

	for i := 0; i < nPayments; i++ {
		paymentArg := CreateReceiptTxPaymentParams{
			PaymentDatetime: arg.ReceiptDatetime,
			Amount:          util.RandomPaymentAmount(),
			PaymentMethodID: paymentMethod.PaymentMethodID,
		}

		arg.ReceiptPaymentsParams = append(arg.ReceiptPaymentsParams, paymentArg)
	}

	result, err := store.CreateReceiptWithPaymentsTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// check Receipt
	require.NotEmpty(t, result.Receipt)
	require.NotZero(t, result.Receipt.ReceiptID)

	receipt, err := testQueries.GetReceipt(context.Background(), result.Receipt.ReceiptID)
	require.NoError(t, err)
	require.NotEmpty(t, receipt)

	require.Equal(t, result.Receipt.ReceiptID, receipt.ReceiptID)
	require.Equal(t, result.Receipt.StudentID, receipt.StudentID)
	require.WithinDuration(t, result.Receipt.ReceiptDatetime, receipt.ReceiptDatetime, time.Second)
	require.Equal(t, result.Receipt.Notes, receipt.Notes)

	// check Payments
	require.NotEmpty(t, result.Payments)
	require.Equal(t, len(result.Payments), nPayments)

	// check each Payment in Payments
	amount := 0.0
	for _, v := range result.Payments {
		require.NotEmpty(t, v)
		require.NotZero(t, v.PaymentID)

		payment, err := testQueries.GetPayment(context.Background(), v.PaymentID)
		require.NoError(t, err)
		require.NotEmpty(t, payment)

		require.Equal(t, v.PaymentID, payment.PaymentID)
		require.Equal(t, v.ReceiptID, payment.ReceiptID)
		require.WithinDuration(t, v.PaymentDatetime, payment.PaymentDatetime, time.Second)
		require.Equal(t, v.Amount, payment.Amount)
		require.Equal(t, v.PaymentMethodID, payment.PaymentMethodID)

		amount += payment.Amount
	}

	require.Equal(t, receipt.Amount, amount)

	return result
}

// createRandomStudentReceiptsWithPaymentsTx adds a single Student with 'n' random ReceiptsWithPayments to the database.
func createRandomStudentReceiptsWithPaymentsTx(t *testing.T, nReceipts int) StudentReceiptsWithPayments {
	var result StudentReceiptsWithPayments
	nPayments := 2 // number of payments created for each receipt

	store := NewStore(testDB)
	student := createRandomStudent(t)
	paymentMethod := createRandomPaymentMethod(t)

	result.StudentID = student.StudentID

	for i := 0; i < nReceipts; i++ {
		arg := CreateReceiptTxParams{
			StudentID:       student.StudentID,
			ReceiptDatetime: util.RandomDatetime(),
			Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
		}

		for j := 0; j < nPayments; j++ {
			paymentArg := CreateReceiptTxPaymentParams{
				PaymentDatetime: arg.ReceiptDatetime,
				Amount:          util.RandomPaymentAmount(),
				PaymentMethodID: paymentMethod.PaymentMethodID,
			}

			arg.ReceiptPaymentsParams = append(arg.ReceiptPaymentsParams, paymentArg)
		}

		receiptWithPayments, err := store.CreateReceiptWithPaymentsTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, receiptWithPayments)

		// check Receipt
		require.NotEmpty(t, receiptWithPayments.Receipt)
		require.NotZero(t, receiptWithPayments.Receipt.ReceiptID)

		receipt, err := testQueries.GetReceipt(context.Background(), receiptWithPayments.Receipt.ReceiptID)
		require.NoError(t, err)
		require.NotEmpty(t, receipt)

		require.Equal(t, receipt.ReceiptID, receiptWithPayments.Receipt.ReceiptID)
		require.Equal(t, receipt.StudentID, receiptWithPayments.Receipt.StudentID)
		require.WithinDuration(t, receipt.ReceiptDatetime, receiptWithPayments.Receipt.ReceiptDatetime, time.Second)
		require.Equal(t, receipt.Notes, receiptWithPayments.Receipt.Notes)

		// check Payments
		require.NotEmpty(t, receiptWithPayments.Payments)
		require.Equal(t, len(receiptWithPayments.Payments), nPayments)

		// check each Payment in Payments
		amount := 0.0
		for _, v := range receiptWithPayments.Payments {
			require.NotEmpty(t, v)
			require.NotZero(t, v.PaymentID)

			payment, err := testQueries.GetPayment(context.Background(), v.PaymentID)
			require.NoError(t, err)
			require.NotEmpty(t, payment)

			require.Equal(t, payment.PaymentID, v.PaymentID)
			require.Equal(t, payment.ReceiptID, v.ReceiptID)
			require.WithinDuration(t, payment.PaymentDatetime, v.PaymentDatetime, time.Second)
			require.Equal(t, payment.Amount, v.Amount)
			require.Equal(t, payment.PaymentMethodID, v.PaymentMethodID)

			amount += payment.Amount
		}

		require.Equal(t, receipt.Amount, amount)
		result.ReceiptsWithPayments = append(result.ReceiptsWithPayments, receiptWithPayments)
	}

	sort.Sort(result.ReceiptsWithPayments)
	return result
}

func TestCreateReceiptWithPaymentsTx(t *testing.T) {
	createRandomReceiptWithPaymentsTx(t, 2)
}

func TestGetReceiptWithPaymentsTx(t *testing.T) {
	store := NewStore(testDB)
	receiptWithPayments1 := createRandomReceiptWithPaymentsTx(t, 5)

	receiptWithPayments2, err := store.GetReceiptWithPaymentsTx(context.Background(), receiptWithPayments1.Receipt.ReceiptID)
	require.NoError(t, err)
	require.NotEmpty(t, receiptWithPayments2)

	//check Receipt
	require.NotEmpty(t, receiptWithPayments2.Receipt)

	receipt1 := receiptWithPayments1.Receipt
	receipt2 := receiptWithPayments2.Receipt

	require.Equal(t, receipt1.ReceiptID, receipt2.ReceiptID)
	require.Equal(t, receipt1.StudentID, receipt2.StudentID)
	require.WithinDuration(t, receipt1.ReceiptDatetime, receipt2.ReceiptDatetime, time.Second)
	require.Equal(t, receipt1.Amount, receipt2.Amount)
	require.Equal(t, receipt1.Notes, receipt2.Notes)

	// check Payments
	require.NotEmpty(t, receiptWithPayments2.Payments)
	require.Equal(t, len(receiptWithPayments2.Payments), len(receiptWithPayments1.Payments))

	// check each Payment in Payments
	for i, payment2 := range receiptWithPayments2.Payments {
		require.NotEmpty(t, payment2)

		payment1 := receiptWithPayments1.Payments[i]

		require.Equal(t, payment1.PaymentID, payment2.PaymentID)
		require.Equal(t, payment1.ReceiptID, payment2.ReceiptID)
		require.WithinDuration(t, payment1.PaymentDatetime, payment2.PaymentDatetime, time.Second)
		require.Equal(t, payment1.Amount, payment2.Amount)
		require.Equal(t, payment1.PaymentMethodID, payment2.PaymentMethodID)
	}
}

func TestDeleteReceiptWithPaymentsTx(t *testing.T) {
	store := NewStore(testDB)

	receiptWithPayments := createRandomReceiptWithPaymentsTx(t, 2)

	err := store.DeleteReceiptWithPaymentsTx(context.Background(), receiptWithPayments.Receipt.ReceiptID)
	require.NoError(t, err)

	// check receipt deleted
	receipt, err := testQueries.GetReceipt(context.Background(), receiptWithPayments.Receipt.ReceiptID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, receipt)

	// check payments deleted
	for _, v := range receiptWithPayments.Payments {
		payment, err := testQueries.GetPayment(context.Background(), v.PaymentID)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, payment)
	}
}

func TestGetReceiptsWithPaymentsByStudentTx(t *testing.T) {
	store := NewStore(testDB)

	nReceipts := 5
	studentReceiptsWithPayments1 := createRandomStudentReceiptsWithPaymentsTx(t, nReceipts)

	studentReceiptsWithPayments2, err := store.GetReceiptsWithPaymentsByStudentTx(context.Background(), studentReceiptsWithPayments1.StudentID, nReceipts+1, 0)
	require.NoError(t, err)
	require.NotEmpty(t, studentReceiptsWithPayments2)
	require.Equal(t, studentReceiptsWithPayments2.StudentID, studentReceiptsWithPayments1.StudentID)

	// check Receipts
	require.NotEmpty(t, studentReceiptsWithPayments2.ReceiptsWithPayments)
	require.Equal(t, len(studentReceiptsWithPayments2.ReceiptsWithPayments), nReceipts)

	// check each ReceiptWithPayments
	for i, receiptWithPayments2 := range studentReceiptsWithPayments2.ReceiptsWithPayments {
		require.NotEmpty(t, receiptWithPayments2)

		// check Receipt
		receiptWithPayments1 := studentReceiptsWithPayments1.ReceiptsWithPayments[i]

		require.Equal(t, receiptWithPayments2.Receipt.ReceiptID, receiptWithPayments1.Receipt.ReceiptID)
		require.Equal(t, receiptWithPayments2.Receipt.StudentID, receiptWithPayments1.Receipt.StudentID)
		require.WithinDuration(t, receiptWithPayments2.Receipt.ReceiptDatetime, receiptWithPayments1.Receipt.ReceiptDatetime, time.Second)
		require.Equal(t, receiptWithPayments2.Receipt.Amount, receiptWithPayments1.Receipt.Amount)
		require.Equal(t, receiptWithPayments2.Receipt.Notes, receiptWithPayments1.Receipt.Notes)

		// check Payments
		require.NotEmpty(t, receiptWithPayments2.Payments)
		require.Equal(t, len(receiptWithPayments2.Payments), len(receiptWithPayments1.Payments))

		// check each Payment in Payments
		for j, payment2 := range receiptWithPayments2.Payments {
			require.NotEmpty(t, payment2)

			payment1 := receiptWithPayments1.Payments[j]

			require.Equal(t, payment2.PaymentID, payment1.PaymentID)
			require.Equal(t, payment2.ReceiptID, payment1.ReceiptID)
			require.WithinDuration(t, payment2.PaymentDatetime, payment1.PaymentDatetime, time.Second)
			require.Equal(t, payment2.Amount, payment1.Amount)
			require.Equal(t, payment2.PaymentMethodID, payment1.PaymentMethodID)
		}
	}
}
