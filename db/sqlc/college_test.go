package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomCollege tests adding a new random college to the database, and returns the College data type.
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
	require.Equal(t, college1.Name, college2.Name)
}

func TestUpdateCollege(t *testing.T) {
	college1 := createRandomCollege(t)

	arg := UpdateCollegeParams{
		CollegeID: college1.CollegeID,
		Name:      util.RandomName(),
	}
	err := testQueries.UpdateCollege(context.Background(), arg)
	require.NoError(t, err)

	college2, err := testQueries.GetCollege(context.Background(), arg.CollegeID)
	require.NoError(t, err)
	require.NotEmpty(t, college2)

	require.Equal(t, college1.CollegeID, college2.CollegeID)
	require.Equal(t, arg.Name, college2.Name)
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
	require.Len(t, colleges, 5)

	for _, college := range colleges {
		require.NotEmpty(t, college)
	}
}
