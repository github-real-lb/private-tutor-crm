package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

func TestCreateLessonTx(t *testing.T) {
	store := NewStore(testDB)
	nLessons := 10 // number of lessons to be created (and number of concurrent transactions to run).
	nInvoices := 2 // number of invoices (and students) to be created for each lesson.

	var students []Student
	for i := 0; i < nInvoices; i++ {
		students = append(students, createRandomStudent(t))
	}

	errs := make(chan error)
	results := make(chan LessonTxResult)

	// run lesson invoicing transactions concurrently.
	for i := 0; i < nLessons; i++ {
		go func() {
			location := createRandomLessonLocation(t)
			subject := createRandomLessonSubject(t)

			arg := CreateLessonTxParams{
				LessonDatetime: util.RandomDatetime(),
				Duration:       util.RandomLessonDuration(),
				LocationID:     location.LocationID,
				SubjectID:      subject.SubjectID,
				Notes:          sql.NullString{String: util.RandomNote(), Valid: true},
			}

			for j := 0; j < nInvoices; j++ {
				invoiceArg := CreateLessonTxInvoiceParams{
					StudentID: students[j].StudentID,
					HourlyFee: students[j].HourlyFee.Float64,
					Duration:  arg.Duration,
					Discount:  util.RandomDiscount(),
					Amount:    students[j].HourlyFee.Float64 * float64(arg.Duration) / 60.0,
					Notes:     sql.NullString{String: util.RandomNote(), Valid: true},
				}

				arg.LessonInvoicesParams = append(arg.LessonInvoicesParams, invoiceArg)
			}

			result, err := store.CreateLessonTx(context.Background(), arg)

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < nLessons; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check Lesson
		require.NotZero(t, result.Lesson.LessonID)
		require.NotZero(t, result.Lesson.LessonID)

		_, err = testQueries.GetLesson(context.Background(), result.Lesson.LessonID)
		require.NoError(t, err)

		// check Invoices
		require.NotEmpty(t, result.Invoices)
		require.Equal(t, len(result.Invoices), nInvoices)

		// check each Invoice in Invoices
		for _, invoice := range result.Invoices {
			require.NotEmpty(t, invoice)
			require.NotZero(t, invoice.InvoiceID)
			require.Equal(t, invoice.LessonID, result.Lesson.LessonID)
			require.True(t, invoice.Duration <= result.Lesson.Duration)

			_, err = testQueries.GetInvoice(context.Background(), invoice.InvoiceID)
			require.NoError(t, err)
		}
	}
}
