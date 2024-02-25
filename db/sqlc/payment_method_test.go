package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomReference adds a new random reference table
// (PaymentMethod, Funnel, LessonLocation, LessonSubject or PaymentMethod) to the database,
// and returns it.
func createRandomPaymentMethod(t *testing.T) PaymentMethod {
	name := util.RandomName()
	paymentMethod, err := testQueries.CreatePaymentMethod(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, paymentMethod)
	require.Equal(t, name, paymentMethod.Name)
	require.NotZero(t, paymentMethod.PaymentMethodID)
	return paymentMethod
}

func TestCreatePaymentMethod(t *testing.T) {
	createRandomPaymentMethod(t)
}

func TestGetPaymentMethod(t *testing.T) {
	paymentMethod1 := createRandomPaymentMethod(t)
	paymentMethod2, err := testQueries.GetPaymentMethod(context.Background(), paymentMethod1.PaymentMethodID)

	require.NoError(t, err)
	require.NotEmpty(t, paymentMethod2)
	require.Equal(t, paymentMethod1.PaymentMethodID, paymentMethod2.PaymentMethodID)
	require.Equal(t, paymentMethod1.PaymentMethodID, paymentMethod2.PaymentMethodID)

}

func TestDeletePaymentMethod(t *testing.T) {
	paymentMethod1 := createRandomPaymentMethod(t)

	err := testQueries.DeletePaymentMethod(context.Background(), paymentMethod1.PaymentMethodID)
	require.NoError(t, err)

	paymentMethod2, err := testQueries.GetPaymentMethod(context.Background(), paymentMethod1.PaymentMethodID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, paymentMethod2)
}

func TestUpdatePaymentMethod(t *testing.T) {
	paymentMethod1 := createRandomPaymentMethod(t)
	name := util.RandomName()

	arg := UpdatePaymentMethodParams{
		PaymentMethodID: paymentMethod1.PaymentMethodID,
		Name:            name,
	}
	err := testQueries.UpdatePaymentMethod(context.Background(), arg)
	require.NoError(t, err)

	paymentMethod2, err := testQueries.GetPaymentMethod(context.Background(), paymentMethod1.PaymentMethodID)
	require.NoError(t, err)
	require.NotEmpty(t, paymentMethod2)

	require.Equal(t, paymentMethod1.PaymentMethodID, paymentMethod2.PaymentMethodID)
	require.Equal(t, name, paymentMethod2.Name)
}

func TestListPaymentMethods(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPaymentMethod(t)
	}

	arg := ListPaymentMethodsParams{
		Limit:  5,
		Offset: 5,
	}
	paymentMethods, err := testQueries.ListPaymentMethods(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, paymentMethods, int(arg.Limit))

	for _, v := range paymentMethods {
		require.NotEmpty(t, v)
	}
}
