package db

import (
	"context"
	"database/sql"
	"time"
)

// CreateLessonTxInvoiceParams contains the input paramaters of a single invoice, for the CreateLessonTx function.
type CreateLessonTxInvoiceParams struct {
	StudentID int64          `json:"student_id"`
	HourlyFee float64        `json:"hourly_fee"`
	Duration  int64          `json:"duration"`
	Discount  float64        `json:"discount"`
	Amount    float64        `json:"amount"`
	Notes     sql.NullString `json:"notes"`
}

// CreateLessonTxParams contains the input paramaters of a single lesson and its Invoices, for the CreateLessonTx function.
type CreateLessonTxParams struct {
	LessonDatetime       time.Time                     `json:"lesson_datetime"`
	Duration             int64                         `json:"duration"`
	LocationID           int64                         `json:"location_id"`
	SubjectID            int64                         `json:"subject_id"`
	Notes                sql.NullString                `json:"notes"`
	LessonInvoicesParams []CreateLessonTxInvoiceParams `json:"lesson_invoices_params"`
}

// LessonTxResult is the result lesson and its invoices, for any lesson transaction.
type LessonTxResult struct {
	Lesson   Lesson    `json:"lesson"`
	Invoices []Invoice `json:"invoices"`
}

// CreateLessonTx creates a lesson that took place,
// and invoices for all the students that took part in the lesson.
func (store *Store) CreateLessonTx(ctx context.Context, arg CreateLessonTxParams) (LessonTxResult, error) {
	var result LessonTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		createLessonArg := CreateLessonParams{
			LessonDatetime: arg.LessonDatetime,
			Duration:       arg.Duration,
			LocationID:     arg.LocationID,
			SubjectID:      arg.SubjectID,
			Notes:          arg.Notes,
		}

		result.Lesson, err = q.CreateLesson(ctx, createLessonArg)
		if err != nil {
			return err
		}

		for _, invoiceArg := range arg.LessonInvoicesParams {
			createInvoiceArg := CreateInvoiceParams{
				StudentID: invoiceArg.StudentID,
				LessonID:  result.Lesson.LessonID,
				HourlyFee: invoiceArg.HourlyFee,
				Duration:  invoiceArg.Duration,
				Discount:  invoiceArg.Discount,
				Amount:    invoiceArg.Amount,
				Notes:     invoiceArg.Notes,
			}

			invoice, err := q.CreateInvoice(ctx, createInvoiceArg)
			if err != nil {
				return err
			}

			result.Invoices = append(result.Invoices, invoice)
		}

		return nil
	})

	return result, err
}
