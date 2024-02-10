package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomInvoice tests adding a new random invoice to the database, and returns the Invoice data type.
func createRandomInvoice(t *testing.T) Invoice {
	student := createRandomStudent(t)
	lesson := createRandomLesson(t)

	arg := CreateInvoiceParams{
		StudentID: student.StudentID,
		LessonID:  lesson.LessonID,
		HourlyFee: util.RandomLessonHourlyFee(),
		Duration:  util.RandomLessonDuration(),
		Discount:  util.RandomDiscount(),
		Amount:    util.RandomInvoiceAmount(),
		Notes:     sql.NullString{String: util.RandomNote(), Valid: true},
	}

	invoice, err := testQueries.CreateInvoice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invoice)

	require.Equal(t, arg.StudentID, invoice.StudentID)
	require.Equal(t, arg.LessonID, invoice.LessonID)
	require.Equal(t, arg.HourlyFee, invoice.HourlyFee)
	require.Equal(t, arg.Duration, invoice.Duration)
	require.Equal(t, arg.Discount, invoice.Discount)
	require.Equal(t, arg.Amount, invoice.Amount)
	require.Equal(t, arg.Notes, invoice.Notes)

	require.NotZero(t, invoice.InvoiceID)
	require.NotZero(t, invoice.CreatedAt)

	return invoice
}

func TestCreateInvoice(t *testing.T) {
	createRandomInvoice(t)
}

func TestGetInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)
	invoice2, err := testQueries.GetInvoice(context.Background(), invoice1.InvoiceID)
	require.NoError(t, err)
	require.NotEmpty(t, invoice2)

	require.Equal(t, invoice1.InvoiceID, invoice2.InvoiceID)
	require.Equal(t, invoice1.StudentID, invoice2.StudentID)
	require.Equal(t, invoice1.LessonID, invoice2.LessonID)
	require.Equal(t, invoice1.HourlyFee, invoice2.HourlyFee)
	require.Equal(t, invoice1.Duration, invoice2.Duration)
	require.Equal(t, invoice1.Discount, invoice2.Discount)
	require.Equal(t, invoice1.Amount, invoice2.Amount)
	require.Equal(t, invoice1.Notes, invoice2.Notes)
	require.WithinDuration(t, invoice1.CreatedAt, invoice2.CreatedAt, time.Second)
}

func TestUpdateInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)
	student := createRandomStudent(t)
	lesson := createRandomLesson(t)

	arg := UpdateInvoiceParams{
		InvoiceID: invoice1.InvoiceID,
		StudentID: student.StudentID,
		LessonID:  lesson.LessonID,
		HourlyFee: util.RandomLessonHourlyFee(),
		Duration:  util.RandomLessonDuration(),
		Discount:  util.RandomDiscount(),
		Amount:    util.RandomInvoiceAmount(),
		Notes:     sql.NullString{String: util.RandomNote(), Valid: true},
	}
	err := testQueries.UpdateInvoice(context.Background(), arg)
	require.NoError(t, err)

	invoice2, err := testQueries.GetInvoice(context.Background(), arg.InvoiceID)
	require.NoError(t, err)
	require.NotEmpty(t, invoice2)

	require.Equal(t, arg.InvoiceID, invoice2.InvoiceID)
	require.Equal(t, arg.StudentID, invoice2.StudentID)
	require.Equal(t, arg.LessonID, invoice2.LessonID)
	require.Equal(t, arg.Duration, invoice2.Duration)
	require.Equal(t, arg.Discount, invoice2.Discount)
	require.Equal(t, arg.HourlyFee, invoice2.HourlyFee)
	require.Equal(t, arg.Amount, invoice2.Amount)
	require.Equal(t, arg.Notes, invoice2.Notes)
	require.WithinDuration(t, invoice1.CreatedAt, invoice2.CreatedAt, time.Second)
}

func TestDeleteInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)
	require.NotEmpty(t, invoice1)
	require.NotZero(t, invoice1.InvoiceID)

	err := testQueries.DeleteInvoice(context.Background(), invoice1.InvoiceID)
	require.NoError(t, err)

	invoice2, err := testQueries.GetInvoice(context.Background(), invoice1.InvoiceID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, invoice2)
}

func TestListInvoices(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomInvoice(t)
	}

	arg := ListInvoicesParams{
		Limit:  5,
		Offset: 5,
	}

	invoices, err := testQueries.ListInvoices(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, invoices, 5)

	for _, invoice := range invoices {
		require.NotEmpty(t, invoice)
	}
}
