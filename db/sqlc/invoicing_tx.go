package db

import (
	"context"
	"database/sql"
	"time"
)

type Invoices []Invoice

// Implement the sort.Interface methods for the Invoices type
func (inv Invoices) Len() int           { return len(inv) }
func (inv Invoices) Swap(i, j int)      { inv[i], inv[j] = inv[j], inv[i] }
func (inv Invoices) Less(i, j int) bool { return inv[i].InvoiceDatetime.Before(inv[j].InvoiceDatetime) }

type Lessons []Lesson

// Implement the sort.Interface methods for the Lessons type
func (l Lessons) Len() int           { return len(l) }
func (l Lessons) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Lessons) Less(i, j int) bool { return l[i].LessonDatetime.Before(l[j].LessonDatetime) }

// LessonWithInvoices is used for a single lesson and all the Invoices issued for participating students.
type LessonWithInvoices struct {
	Lesson   Lesson   `json:"lesson"`
	Invoices Invoices `json:"invoices"`
}

// CreateLessonTxInvoiceParams contains the input paramaters of a single invoice, for the CreateLessonWithInvoicesTx function.
type CreateLessonTxInvoiceParams struct {
	StudentID int64          `json:"student_id"`
	HourlyFee float64        `json:"hourly_fee"`
	Duration  int64          `json:"duration"`
	Discount  float64        `json:"discount"`
	Amount    float64        `json:"amount"`
	Notes     sql.NullString `json:"notes"`
}

// CreateLessonTxParams contains the input paramaters of a single lesson and its Invoices, for the CreateLessonWithInvoicesTx function.
type CreateLessonTxParams struct {
	LessonDatetime       time.Time                     `json:"lesson_datetime"`
	Duration             int64                         `json:"duration"`
	LocationID           int64                         `json:"location_id"`
	SubjectID            int64                         `json:"subject_id"`
	Notes                sql.NullString                `json:"notes"`
	LessonInvoicesParams []CreateLessonTxInvoiceParams `json:"lesson_invoices_params"`
}

// CreateLessonTx creates a lesson held and invoices for all the students that took part in the lesson.
func (store *SQLStore) CreateLessonWithInvoicesTx(ctx context.Context, arg CreateLessonTxParams) (LessonWithInvoices, error) {
	var result LessonWithInvoices

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

// GetLessonWithInvoicesTx gets a Lesson and all the Invoices releated to it.
func (store *SQLStore) GetLessonWithInvoicesTx(ctx context.Context, lessonID int64) (LessonWithInvoices, error) {
	var result LessonWithInvoices

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Lesson, err = q.GetLesson(ctx, lessonID)
		if err != nil {
			return err
		}

		result.Invoices, err = q.GetInvoicesByLesson(ctx, lessonID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// DeleteLessonWithInvoicesTx deletes a Lesson and all the Invoices releated to it.
func (store *SQLStore) DeleteLessonWithInvoicesTx(ctx context.Context, lessonID int64) error {
	err := store.execTx(ctx, func(q *Queries) error {
		err := q.DeleteInvoicesByLesson(ctx, lessonID)
		if err != nil {
			return err
		}

		err = q.DeleteLesson(ctx, lessonID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
