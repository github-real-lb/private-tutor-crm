package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

func TestLessonInvoicingTx(t *testing.T) {
	store := NewStore(testDB)
	nLessons := 10 // number of lessons to be created (and number of concurrent transactions to run).
	nInvoices := 2 // number of invoices (and students) to be created for each lesson.

	var students []Student
	for i := 0; i < nInvoices; i++ {
		students = append(students, createRandomStudent(t))
	}

	errs := make(chan error)
	results := make(chan LessonInvoicingTxResult)

	// run lesson invoicing transactions concurrently.
	for i := 0; i < nLessons; i++ {
		go func() {
			location := createRandomLessonLocation(t)
			subject := createRandomLessonSubject(t)

			arg := LessonInvoicingTxParams{
				LessonDatetime: util.RandomDatetime(),
				Duration:       util.RandomLessonDuration(),
				LocationID:     location.LocationID,
				SubjectID:      subject.SubjectID,
				Notes:          sql.NullString{String: util.RandomNote(), Valid: true},
			}

			for j := 0; j < nInvoices; j++ {
				invoiceArg := LessonInvoicingTxInvoiceParams{
					StudentID: students[j].StudentID,
					HourlyFee: students[j].HourlyFee.Float64,
					Duration:  arg.Duration,
					Discount:  util.RandomDiscount(),
					Amount:    students[j].HourlyFee.Float64 * float64(arg.Duration) / 60.0,
					Notes:     sql.NullString{String: util.RandomNote(), Valid: true},
				}

				arg.LessonInvoicesParams = append(arg.LessonInvoicesParams, invoiceArg)
			}

			result, err := store.LessonInvoicingTx(context.Background(), arg)

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < nLessons; i++ {
		err := <-errs
		require.NoError(t, err)

		//check Lesson
		result := <-results
		require.NotEmpty(t, result)

		_, err = testQueries.GetLesson(context.Background(), result.Lesson.LessonID)
		require.NoError(t, err)

		// check Invoices
		require.NotEmpty(t, result.Invoices)
		require.Equal(t, len(result.Invoices), nInvoices)

		// check each Invoice in Invoices
		for _, invoice := range result.Invoices {
			require.NotEmpty(t, invoice)

			_, err = testQueries.GetInvoice(context.Background(), invoice.InvoiceID)
			require.NoError(t, err)
			require.Equal(t, invoice.LessonID, result.Lesson.LessonID)
		}
	}
}

func TestPaymentsReceivingTx(t *testing.T) {
	store := NewStore(testDB)
	nRecipts := 10 // Number of receipts (and students) to be created (and number of concurrent transactions to run).
	nPayments := 2 // Number of payments to be created for each recipt

	var students []Student
	for i := 0; i < nRecipts; i++ {
		students = append(students, createRandomStudent(t))
	}

	var paymentMethods []PaymentMethod
	for i := 0; i < nPayments; i++ {
		paymentMethods = append(paymentMethods, createRandomPaymentMethod(t))
	}

	n := make(chan int) // channel of recipt transaction number
	errs := make(chan error)
	results := make(chan PaymentsReceivingTxResult)

	// run payment receiving transactions concurrently.
	for i := 0; i < nRecipts; i++ {
		go func() {
			arg := PaymentsReceivingTxParams{
				StudentID:       students[<-n].StudentID,
				ReceiptDatetime: util.RandomDatetime(),
				Amount:          0.0,
				Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
			}

			for j := 0; j < nPayments; j++ {
				paymentArg := PaymentsReceivingTxPaymentParams{
					PaymentDatetime: arg.ReceiptDatetime,
					Amount:          util.RandomPaymentAmount(),
					PaymentMethodID: paymentMethods[j].PaymentMethodID,
				}

				arg.Amount += paymentArg.Amount
				arg.ReceiptPaymentsParams = append(arg.ReceiptPaymentsParams, paymentArg)
			}

			result, err := store.PaymentsReceivingTx(context.Background(), arg)

			errs <- err
			results <- result
		}()

		n <- i
	}

	// check results
	for i := 0; i < nRecipts; i++ {
		err := <-errs
		require.NoError(t, err)

		// check Receipt
		result := <-results
		require.NotEmpty(t, result)

		_, err = testQueries.GetReceipt(context.Background(), result.Receipt.ReceiptID)
		require.NoError(t, err)

		// check Payments
		require.NotEmpty(t, result.Payments)
		require.Equal(t, len(result.Payments), nPayments)

		// check each Payment in Payments
		for _, payment := range result.Payments {
			require.NotEmpty(t, payment)

			_, err = testQueries.GetPayment(context.Background(), payment.PaymentID)
			require.NoError(t, err)
			require.Equal(t, payment.ReceiptID, result.Receipt.ReceiptID)
		}
	}
}
