package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomReference adds a new random reference table
// (College, Funnel, LessonLocation, LessonSubject or PaymentMethod) to the database,
// and returns it.
func createRandomReference(t *testing.T, key string) Reference {
	ref, ok := ReferencesMap[key]
	require.True(t, ok)

	name := util.RandomName()

	switch key {
	case "College":
		college, err := testQueries.CreateCollege(context.Background(), name)
		require.NoError(t, err)
		ref = &college
	case "Funnel":
		funnel, err := testQueries.CreateFunnel(context.Background(), name)
		require.NoError(t, err)
		ref = &funnel
	case "LessonLocation":
		location, err := testQueries.CreateLessonLocation(context.Background(), name)
		require.NoError(t, err)
		ref = &location
	case "LessonSubject":
		subject, err := testQueries.CreateLessonSubject(context.Background(), name)
		require.NoError(t, err)
		ref = &subject
	case "PaymentMethod":
		paymentMethod, err := testQueries.CreatePaymentMethod(context.Background(), name)
		require.NoError(t, err)
		ref = &paymentMethod
	}

	require.NotEmpty(t, ref)
	require.Equal(t, name, ref.GetName())
	require.NotZero(t, ref.GetID())
	return ref

}

func TestCreateReferences(t *testing.T) {
	for key := range ReferencesMap {
		t.Run(key, func(t *testing.T) {
			createRandomReference(t, key)
		})
	}
}

func TestGetReferences(t *testing.T) {
	for key := range ReferencesMap {
		t.Run(key, func(t *testing.T) {
			ref1 := createRandomReference(t, key)

			var ref2 Reference
			switch key {
			case "College":
				college, err := testQueries.GetCollege(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &college
			case "Funnel":
				funnel, err := testQueries.GetFunnel(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &funnel
			case "LessonLocation":
				location, err := testQueries.GetLessonLocation(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &location
			case "LessonSubject":
				subject, err := testQueries.GetLessonSubject(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &subject
			case "PaymentMethod":
				paymentMethod, err := testQueries.GetPaymentMethod(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &paymentMethod
			}

			require.NotEmpty(t, ref2)
			require.Equal(t, ref1.GetID(), ref2.GetID())
			require.Equal(t, ref1.GetName(), ref2.GetName())

		})
	}
}

func TestDeleteReferences(t *testing.T) {
	for key := range ReferencesMap {
		t.Run(key, func(t *testing.T) {
			ref1 := createRandomReference(t, key)

			var ref2 Reference
			switch key {
			case "College":
				err := testQueries.DeleteCollege(context.Background(), ref1.GetID())
				require.NoError(t, err)
				college, err := testQueries.GetCollege(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &college
			case "Funnel":
				funnel, err := testQueries.GetFunnel(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &funnel
			case "LessonLocation":
				location, err := testQueries.GetLessonLocation(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &location
			case "LessonSubject":
				subject, err := testQueries.GetLessonSubject(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &subject
			case "PaymentMethod":
				paymentMethod, err := testQueries.GetPaymentMethod(context.Background(), ref1.GetID())
				require.NoError(t, err)
				ref2 = &paymentMethod
			}

			require.NotEmpty(t, ref2)
			require.Equal(t, ref1.GetID(), ref2.GetID())
			require.Equal(t, ref1.GetName(), ref2.GetName())

		})
	}
	college1 := createRandomCollege(t)

	err := testQueries.DeleteCollege(context.Background(), college1.CollegeID)
	require.NoError(t, err)

	college2, err := testQueries.GetCollege(context.Background(), college1.CollegeID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, college2)
}
