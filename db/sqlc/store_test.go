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
	nLessons := 5  // number of lessons to be created
	nInvoices := 2 // number of invoices to be created

	var students []Student
	for i := 0; i < nLessons+1; i++ {
		students = append(students, createRandomStudent(t))
	}

	for i := 0; i < nLessons; i++ {
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
				StudentID: students[i+j].StudentID,
				HourlyFee: students[i+j].HourlyFee.Float64,
				Duration:  arg.Duration,
				Discount:  util.RandomDiscount(),
				Amount:    students[i+j].HourlyFee.Float64 * float64(arg.Duration) / 60.0,
				Notes:     sql.NullString{String: util.RandomNote(), Valid: true},
			}

			arg.LessonInvoicesParams = append(arg.LessonInvoicesParams, invoiceArg)
		}

		result, err := store.LessonInvoicingTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Invoices)
	}
}

func TestPaymentsReceivingTx(t *testing.T) {
	store := NewStore(testDB)
	nRecipts := 5        // Number of receipts to be created
	nPaymentMethods := 2 // Number of payment methods to be created

	var students []Student
	for i := 0; i < nRecipts; i++ {
		students = append(students, createRandomStudent(t))
	}

	var paymentMethods []PaymentMethod
	for i := 0; i < nPaymentMethods; i++ {
		paymentMethods = append(paymentMethods, createRandomPaymentMethod(t))
	}

	for i := 0; i < nRecipts; i++ {
		arg := PaymentsReceivingTxParams{
			StudentID:       students[i].StudentID,
			ReceiptDatetime: util.RandomDatetime(),
			Amount:          0.0,
			Notes:           sql.NullString{String: util.RandomNote(), Valid: true},
		}

		for j := 0; j < nPaymentMethods; j++ {
			paymentArg := PaymentsReceivingTxPaymentParams{
				PaymentDatetime: arg.ReceiptDatetime,
				Amount:          util.RandomPaymentAmount(),
				PaymentMethodID: paymentMethods[j].PaymentMethodID,
			}

			arg.Amount += paymentArg.Amount
			arg.ReceiptPaymentsParams = append(arg.ReceiptPaymentsParams, paymentArg)
		}

		result, err := store.PaymentsReceivingTx(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Payments)
	}
}
