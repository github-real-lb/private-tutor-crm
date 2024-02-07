package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomLesson tests adding a new random lesson to the database, and returns the Lesson data type.
func createRandomLesson(t *testing.T) Lesson {
	lessonLocation := createRandomLessonLocation(t)
	lessonSubject := createRandomLessonSubject(t)

	arg := CreateLessonParams{
		LessonDatetime: util.RandomDatetime(),
		Duration:       util.RandomLessonDuration(),
		LocationID:     lessonLocation.LocationID,
		SubjectID:      lessonSubject.SubjectID,
		Notes:          sql.NullString{String: util.RandomNote(), Valid: true},
	}

	lesson, err := testQueries.CreateLesson(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, lesson)

	require.WithinDuration(t, arg.LessonDatetime, lesson.LessonDatetime, time.Second)
	require.Equal(t, arg.Duration, lesson.Duration)
	require.Equal(t, arg.LocationID, lesson.LocationID)
	require.Equal(t, arg.SubjectID, lesson.SubjectID)
	require.Equal(t, arg.Notes, lesson.Notes)

	require.NotZero(t, lesson.LessonID)

	return lesson
}

func TestCreateLesson(t *testing.T) {
	createRandomLesson(t)
}

func TestGetLesson(t *testing.T) {
	lesson1 := createRandomLesson(t)
	lesson2, err := testQueries.GetLesson(context.Background(), lesson1.LessonID)
	require.NoError(t, err)
	require.NotEmpty(t, lesson2)

	require.Equal(t, lesson1.LessonID, lesson2.LessonID)
	require.WithinDuration(t, lesson1.LessonDatetime, lesson2.LessonDatetime, time.Second)
	require.Equal(t, lesson1.Duration, lesson2.Duration)
	require.Equal(t, lesson1.LocationID, lesson2.LocationID)
	require.Equal(t, lesson1.SubjectID, lesson2.SubjectID)
	require.Equal(t, lesson1.Notes, lesson2.Notes)
}

func TestUpdateLesson(t *testing.T) {
	lesson1 := createRandomLesson(t)
	lessonLocation := createRandomLessonLocation(t)
	lessonSubject := createRandomLessonSubject(t)

	arg := UpdateLessonParams{
		LessonID:       lesson1.LessonID,
		LessonDatetime: util.RandomDatetime(),
		Duration:       util.RandomLessonDuration(),
		LocationID:     lessonLocation.LocationID,
		SubjectID:      lessonSubject.SubjectID,
		Notes:          sql.NullString{String: util.RandomNote(), Valid: true},
	}
	err := testQueries.UpdateLesson(context.Background(), arg)
	require.NoError(t, err)

	lesson2, err := testQueries.GetLesson(context.Background(), arg.LessonID)
	require.NoError(t, err)
	require.NotEmpty(t, lesson2)

	require.Equal(t, arg.LessonID, lesson2.LessonID)
	require.WithinDuration(t, arg.LessonDatetime, lesson2.LessonDatetime, time.Second)
	require.Equal(t, arg.Duration, lesson2.Duration)
	require.Equal(t, arg.LocationID, lesson2.LocationID)
	require.Equal(t, arg.SubjectID, lesson2.SubjectID)
	require.Equal(t, arg.Notes, lesson2.Notes)
}

func TestDeleteLesson(t *testing.T) {
	lesson1 := createRandomLesson(t)
	testQueries.DeleteLesson(context.Background(), lesson1.LessonID)

	lesson2, err := testQueries.GetLesson(context.Background(), lesson1.LessonID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, lesson2)
}

func TestListLessons(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomLesson(t)
	}

	arg := ListLessonsParams{
		Limit:  5,
		Offset: 5,
	}

	lessons, err := testQueries.ListLessons(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, lessons, 5)

	for _, lesson := range lessons {
		require.NotEmpty(t, lesson)
	}
}
