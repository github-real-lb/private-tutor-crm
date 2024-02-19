package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomReceipt adds a new random receipt to the database, and returns the Receipt data type.
func createRandomReceipt(t *testing.T) *Receipt {
	student := createRandomStudent(t)

	arg := CreateReceiptParams{
		StudentID:       student.StudentID,
		ReceiptDatetime: util.RandomDatetime(),
		Amount:          util.RandomPaymentAmount(),
		Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
	}

	receipt, err := testQueries.CreateReceipt(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, receipt)

	require.NotZero(t, receipt.ReceiptID)

	require.Equal(t, arg.StudentID, receipt.StudentID)
	require.WithinDuration(t, arg.ReceiptDatetime, receipt.ReceiptDatetime, time.Second)
	require.Equal(t, arg.Amount, receipt.Amount)
	require.Equal(t, arg.Notes, receipt.Notes)

	return receipt
}

func TestCreateReceipt(t *testing.T) {
	createRandomReceipt(t)
}

func TestGetReceipt(t *testing.T) {
	receipt1 := createRandomReceipt(t)
	receipt2, err := testQueries.GetReceipt(context.Background(), receipt1.ReceiptID)
	require.NoError(t, err)
	require.NotEmpty(t, receipt2)

	require.Equal(t, receipt1.ReceiptID, receipt2.ReceiptID)
	require.Equal(t, receipt1.StudentID, receipt2.StudentID)
	require.WithinDuration(t, receipt1.ReceiptDatetime, receipt2.ReceiptDatetime, time.Second)
	require.Equal(t, receipt1.Amount, receipt2.Amount)
	require.Equal(t, receipt1.Notes, receipt2.Notes)
}

func TestUpdateReceipt(t *testing.T) {
	receipt1 := createRandomReceipt(t)
	student := createRandomStudent(t)

	arg := UpdateReceiptParams{
		ReceiptID:       receipt1.ReceiptID,
		StudentID:       student.StudentID,
		ReceiptDatetime: util.RandomDatetime(),
		Amount:          util.RandomPaymentAmount(),
		Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
	}
	err := testQueries.UpdateReceipt(context.Background(), arg)
	require.NoError(t, err)

	receipt2, err := testQueries.GetReceipt(context.Background(), arg.ReceiptID)
	require.NoError(t, err)
	require.NotEmpty(t, receipt2)

	require.Equal(t, arg.ReceiptID, receipt2.ReceiptID)
	require.Equal(t, arg.StudentID, receipt2.StudentID)
	require.WithinDuration(t, arg.ReceiptDatetime, receipt2.ReceiptDatetime, time.Second)
	require.Equal(t, arg.Amount, receipt2.Amount)
	require.Equal(t, arg.Notes, receipt2.Notes)
}

func TestUpdateReceiptAmount(t *testing.T) {
	receipt1 := createRandomReceipt(t)

	arg := UpdateReceiptAmountParams{
		ReceiptID: receipt1.ReceiptID,
		Amount:    util.RandomPaymentAmount(),
	}
	err := testQueries.UpdateReceiptAmount(context.Background(), arg)
	require.NoError(t, err)

	receipt2, err := testQueries.GetReceipt(context.Background(), arg.ReceiptID)
	require.NoError(t, err)
	require.NotEmpty(t, receipt2)

	require.Equal(t, arg.ReceiptID, receipt2.ReceiptID)
	require.Equal(t, arg.Amount, receipt2.Amount)
}

func TestDeleteReceipt(t *testing.T) {
	receipt1 := createRandomReceipt(t)

	err := testQueries.DeleteReceipt(context.Background(), receipt1.ReceiptID)
	require.NoError(t, err)

	receipt2, err := testQueries.GetReceipt(context.Background(), receipt1.ReceiptID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, receipt2)
}

func TestListReceipts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomReceipt(t)
	}

	arg := ListReceiptsParams{
		Limit:  5,
		Offset: 5,
	}

	receipts, err := testQueries.ListReceipts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, receipts, 5)

	for _, receipt := range receipts {
		require.NotEmpty(t, receipt)
	}
}
