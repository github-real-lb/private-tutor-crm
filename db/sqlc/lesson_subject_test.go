package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomLessonSubject tests adding a new random lesson subject to the database, and returns the LessonSubject data type.
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
	require.Equal(t, lessonSubject1.Name, lessonSubject2.Name)
}

func TestUpdateLessonSubject(t *testing.T) {
	lessonSubject1 := createRandomLessonSubject(t)

	arg := UpdateLessonSubjectParams{
		SubjectID: lessonSubject1.SubjectID,
		Name:      util.RandomName(),
	}
	err := testQueries.UpdateLessonSubject(context.Background(), arg)
	require.NoError(t, err)

	lessonSubject2, err := testQueries.GetLessonSubject(context.Background(), arg.SubjectID)
	require.NoError(t, err)
	require.NotEmpty(t, lessonSubject2)

	require.Equal(t, lessonSubject1.SubjectID, lessonSubject2.SubjectID)
	require.Equal(t, arg.Name, lessonSubject2.Name)
}

func TestDeleteLessonSubject(t *testing.T) {
	lessonSubject1 := createRandomLessonSubject(t)
	testQueries.DeleteLessonSubject(context.Background(), lessonSubject1.SubjectID)

	lessonSubject2, err := testQueries.GetLessonSubject(context.Background(), lessonSubject1.SubjectID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, lessonSubject2)
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
	require.Len(t, lessonSubjects, 5)

	for _, lessonSubject := range lessonSubjects {
		require.NotEmpty(t, lessonSubject)
	}
}
