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
func createRandomCollege(t *testing.T) College {
	name := util.RandomName()
	college, err := testQueries.CreateCollege(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, college)
	require.Equal(t, name, college.Name)
	require.NotZero(t, college.CollegeID)
	return college
}

func TestCreateCollege(t *testing.T) {
	createRandomCollege(t)
}

func TestGetCollege(t *testing.T) {
	college1 := createRandomCollege(t)
	college2, err := testQueries.GetCollege(context.Background(), college1.CollegeID)

	require.NoError(t, err)
	require.NotEmpty(t, college2)
	require.Equal(t, college1.CollegeID, college2.CollegeID)
	require.Equal(t, college1.CollegeID, college2.CollegeID)

}

func TestDeleteCollege(t *testing.T) {
	college1 := createRandomCollege(t)

	err := testQueries.DeleteCollege(context.Background(), college1.CollegeID)
	require.NoError(t, err)

	college2, err := testQueries.GetCollege(context.Background(), college1.CollegeID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, college2)
}

func TestUpdateCollege(t *testing.T) {
	college1 := createRandomCollege(t)
	name := util.RandomName()

	arg := UpdateCollegeParams{
		CollegeID: college1.CollegeID,
		Name:      name,
	}
	err := testQueries.UpdateCollege(context.Background(), arg)
	require.NoError(t, err)

	college2, err := testQueries.GetCollege(context.Background(), college1.CollegeID)
	require.NoError(t, err)
	require.NotEmpty(t, college2)

	require.Equal(t, college1.CollegeID, college2.CollegeID)
	require.Equal(t, name, college2.Name)
}

func TestListColleges(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCollege(t)
	}

	arg := ListCollegesParams{
		Limit:  5,
		Offset: 5,
	}
	colleges, err := testQueries.ListColleges(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, colleges, int(arg.Limit))

	for _, v := range colleges {
		require.NotEmpty(t, v)
	}
}
