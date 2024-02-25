package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomReference adds a new random reference table
// (LessonSubject, Funnel, LessonLocation, LessonSubject or PaymentMethod) to the database,
// and returns it.
func createRandomLessonSubject(t *testing.T) LessonSubject {
	name := util.RandomName()
	lessonSubject, err := testQueries.CreateLessonSubject(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, lessonSubject)
	require.Equal(t, name, lessonSubject.Name)
	require.NotZero(t, lessonSubject.SubjectID)
	return lessonSubject
}

func TestCreateLessonSubject(t *testing.T) {
	createRandomLessonSubject(t)
}

func TestGetLessonSubject(t *testing.T) {
	lessonSubject1 := createRandomLessonSubject(t)
	lessonSubject2, err := testQueries.GetLessonSubject(context.Background(), lessonSubject1.SubjectID)

	require.NoError(t, err)
	require.NotEmpty(t, lessonSubject2)
	require.Equal(t, lessonSubject1.SubjectID, lessonSubject2.SubjectID)
	require.Equal(t, lessonSubject1.SubjectID, lessonSubject2.SubjectID)

}

func TestDeleteLessonSubject(t *testing.T) {
	lessonSubject1 := createRandomLessonSubject(t)

	err := testQueries.DeleteLessonSubject(context.Background(), lessonSubject1.SubjectID)
	require.NoError(t, err)

	lessonSubject2, err := testQueries.GetLessonSubject(context.Background(), lessonSubject1.SubjectID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, lessonSubject2)
}

func TestUpdateLessonSubject(t *testing.T) {
	lessonSubject1 := createRandomLessonSubject(t)
	name := util.RandomName()

	arg := UpdateLessonSubjectParams{
		SubjectID: lessonSubject1.SubjectID,
		Name:      name,
	}
	err := testQueries.UpdateLessonSubject(context.Background(), arg)
	require.NoError(t, err)

	lessonSubject2, err := testQueries.GetLessonSubject(context.Background(), lessonSubject1.SubjectID)
	require.NoError(t, err)
	require.NotEmpty(t, lessonSubject2)

	require.Equal(t, lessonSubject1.SubjectID, lessonSubject2.SubjectID)
	require.Equal(t, name, lessonSubject2.Name)
}

func TestListLessonSubjects(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomLessonSubject(t)
	}

	arg := ListLessonSubjectsParams{
		Limit:  5,
		Offset: 5,
	}
	lessonSubjects, err := testQueries.ListLessonSubjects(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, lessonSubjects, int(arg.Limit))

	for _, v := range lessonSubjects {
		require.NotEmpty(t, v)
	}
}
