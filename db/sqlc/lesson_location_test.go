package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomReference adds a new random reference table
// (LessonLocation, Funnel, LessonLocation, LessonSubject or PaymentMethod) to the database,
// and returns it.
func createRandomLessonLocation(t *testing.T) LessonLocation {
	name := util.RandomName()
	lessonLocation, err := testQueries.CreateLessonLocation(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, lessonLocation)
	require.Equal(t, name, lessonLocation.Name)
	require.NotZero(t, lessonLocation.LocationID)
	return lessonLocation
}

func TestCreateLessonLocation(t *testing.T) {
	createRandomLessonLocation(t)
}

func TestGetLessonLocation(t *testing.T) {
	lessonLocation1 := createRandomLessonLocation(t)
	lessonLocation2, err := testQueries.GetLessonLocation(context.Background(), lessonLocation1.LocationID)

	require.NoError(t, err)
	require.NotEmpty(t, lessonLocation2)
	require.Equal(t, lessonLocation1.LocationID, lessonLocation2.LocationID)
	require.Equal(t, lessonLocation1.LocationID, lessonLocation2.LocationID)

}

func TestDeleteLessonLocation(t *testing.T) {
	lessonLocation1 := createRandomLessonLocation(t)

	err := testQueries.DeleteLessonLocation(context.Background(), lessonLocation1.LocationID)
	require.NoError(t, err)

	lessonLocation2, err := testQueries.GetLessonLocation(context.Background(), lessonLocation1.LocationID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, lessonLocation2)
}

func TestUpdateLessonLocation(t *testing.T) {
	lessonLocation1 := createRandomLessonLocation(t)
	name := util.RandomName()

	arg := UpdateLessonLocationParams{
		LocationID: lessonLocation1.LocationID,
		Name:       name,
	}
	err := testQueries.UpdateLessonLocation(context.Background(), arg)
	require.NoError(t, err)

	lessonLocation2, err := testQueries.GetLessonLocation(context.Background(), lessonLocation1.LocationID)
	require.NoError(t, err)
	require.NotEmpty(t, lessonLocation2)

	require.Equal(t, lessonLocation1.LocationID, lessonLocation2.LocationID)
	require.Equal(t, name, lessonLocation2.Name)
}

func TestListLessonLocations(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomLessonLocation(t)
	}

	arg := ListLessonLocationsParams{
		Limit:  5,
		Offset: 5,
	}
	lessonLocations, err := testQueries.ListLessonLocations(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, lessonLocations, int(arg.Limit))

	for _, v := range lessonLocations {
		require.NotEmpty(t, v)
	}
}
