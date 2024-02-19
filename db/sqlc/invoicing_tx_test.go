package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomLessonWithInvoicesTx adds a new random Lesson with 'n' Invoices to the database.
// Each Invoice is created for a different Student participating in the lesson.
func createRandomLessonWithInvoicesTx(t *testing.T, n int) LessonWithInvoices {
	store := NewStore(testDB)

	// create 'n' students
	var students []*Student
	for i := 0; i < n; i++ {
		students = append(students, createRandomStudent(t))
	}

	// create random location and subject for lesson
	location := createRandomReferenceStruct(t, ReferenceLessonLocation)
	subject := createRandomReferenceStruct(t, ReferenceLessonSubject)

	arg := CreateLessonTxParams{
		LessonDatetime: util.RandomDatetime(),
		Duration:       util.RandomLessonDuration(),
		LocationID:     location.GetID(),
		SubjectID:      subject.GetID(),
		Notes:          sql.NullString{String: util.RandomNote(), Valid: true},
	}

	// create 'n' invoices for each student
	for i := 0; i < n; i++ {
		invoiceArg := CreateLessonTxInvoiceParams{
			StudentID: students[i].StudentID,
			HourlyFee: students[i].HourlyFee.Float64,
			Duration:  arg.Duration,
			Discount:  util.RandomDiscount(),
			Amount:    students[i].HourlyFee.Float64 * float64(arg.Duration) / 60.0,
			Notes:     sql.NullString{String: util.RandomNote(), Valid: true},
		}

		// apply the random discount in Discount to calculate Amount
		invoiceArg.Amount *= 1.0 - invoiceArg.Discount

		arg.LessonInvoicesParams = append(arg.LessonInvoicesParams, invoiceArg)
	}

	// create lesson with invoices
	result, err := store.CreateLessonWithInvoicesTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// check Lesson
	require.NotEmpty(t, result.Lesson)
	require.NotZero(t, result.Lesson.LessonID)

	lesson, err := testQueries.GetLesson(context.Background(), result.Lesson.LessonID)
	require.NoError(t, err)

	require.Equal(t, result.Lesson.LessonID, lesson.LessonID)
	require.WithinDuration(t, result.Lesson.LessonDatetime, lesson.LessonDatetime, time.Second)
	require.Equal(t, result.Lesson.Duration, lesson.Duration)
	require.Equal(t, result.Lesson.LocationID, lesson.LocationID)
	require.Equal(t, result.Lesson.SubjectID, lesson.SubjectID)
	require.Equal(t, result.Lesson.Notes, lesson.Notes)

	// check Invoices
	require.NotEmpty(t, result.Invoices)
	require.Equal(t, len(result.Invoices), n)

	// check each Invoice in Invoices
	for _, v := range result.Invoices {
		require.NotEmpty(t, v)
		require.NotZero(t, v.InvoiceID)

		invoice, err := testQueries.GetInvoice(context.Background(), v.InvoiceID)
		require.NoError(t, err)
		require.NotEmpty(t, invoice)

		require.Equal(t, v.InvoiceID, invoice.InvoiceID)
		require.Equal(t, v.StudentID, invoice.StudentID)
		require.Equal(t, v.LessonID, invoice.LessonID)
		require.WithinDuration(t, v.InvoiceDatetime, invoice.InvoiceDatetime, time.Second)
		require.Equal(t, v.HourlyFee, invoice.HourlyFee)
		require.Equal(t, v.Duration, invoice.Duration)
		require.Equal(t, v.Discount, invoice.Discount)
		require.Equal(t, v.Amount, invoice.Amount)
		require.Equal(t, v.Notes, invoice.Notes)

		require.Equal(t, result.Lesson.LessonID, invoice.LessonID)
		require.True(t, result.Lesson.Duration >= invoice.Duration)
	}

	return result
}

func TestCreateLessonWithInvoicesTx(t *testing.T) {
	createRandomLessonWithInvoicesTx(t, 5)
}

func TestGetLessonWithInvoicesTx(t *testing.T) {
	store := NewStore(testDB)
	lessonWithInvoices1 := createRandomLessonWithInvoicesTx(t, 5)

	lessonWithInvoices2, err := store.GetLessonWithInvoicesTx(context.Background(), lessonWithInvoices1.Lesson.LessonID)
	require.NoError(t, err)
	require.NotEmpty(t, lessonWithInvoices2)

	//check Lesson
	require.NotEmpty(t, lessonWithInvoices2.Lesson)

	lesson1 := lessonWithInvoices1.Lesson
	lesson2 := lessonWithInvoices2.Lesson

	require.Equal(t, lesson1.LessonID, lesson2.LessonID)
	require.WithinDuration(t, lesson1.LessonDatetime, lesson2.LessonDatetime, time.Second)
	require.Equal(t, lesson1.Duration, lesson2.Duration)
	require.Equal(t, lesson1.LocationID, lesson2.LocationID)
	require.Equal(t, lesson1.SubjectID, lesson2.SubjectID)
	require.Equal(t, lesson1.Notes, lesson2.Notes)

	// check Invoices
	require.NotEmpty(t, lessonWithInvoices2.Invoices)
	require.Equal(t, len(lessonWithInvoices2.Invoices), len(lessonWithInvoices1.Invoices))

	// check each Invoice in Invoices
	for i, invoice2 := range lessonWithInvoices2.Invoices {
		require.NotEmpty(t, invoice2)

		invoice1 := lessonWithInvoices1.Invoices[i]

		require.Equal(t, invoice1.InvoiceID, invoice2.InvoiceID)
		require.Equal(t, invoice1.StudentID, invoice2.StudentID)
		require.Equal(t, invoice1.LessonID, invoice2.LessonID)
		require.WithinDuration(t, invoice1.InvoiceDatetime, invoice2.InvoiceDatetime, time.Second)
		require.Equal(t, invoice1.HourlyFee, invoice2.HourlyFee)
		require.Equal(t, invoice1.Duration, invoice2.Duration)
		require.Equal(t, invoice1.Discount, invoice2.Discount)
		require.Equal(t, invoice1.Amount, invoice2.Amount)
		require.Equal(t, invoice1.Notes, invoice2.Notes)
	}
}

func TestDeleteLessonWithInvoicesTx(t *testing.T) {
	store := NewStore(testDB)

	LessonWithInvoices := createRandomLessonWithInvoicesTx(t, 5)

	err := store.DeleteLessonWithInvoicesTx(context.Background(), LessonWithInvoices.Lesson.LessonID)
	require.NoError(t, err)

	// check lesson deleted
	lesson, err := testQueries.GetLesson(context.Background(), LessonWithInvoices.Lesson.LessonID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, lesson)

	// check invoices deleted
	for _, v := range LessonWithInvoices.Invoices {
		invoice, err := testQueries.GetInvoice(context.Background(), v.InvoiceID)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, invoice)
	}
}
