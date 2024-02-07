package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomStudent tests adding a new random student to the database, and returns the Student data type.
func createRandomStudent(t *testing.T) Student {
	college, err := testQueries.CreateCollege(context.Background(), util.RandomName())
	require.NoError(t, err)
	require.NotEmpty(t, college)

	funnel, err := testQueries.CreateFunnel(context.Background(), util.RandomName())
	require.NoError(t, err)
	require.NotEmpty(t, funnel)

	arg := CreateStudentParams{
		FirstName:   util.RandomName(),
		LastName:    util.RandomName(),
		Email:       sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber: sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		Address:     sql.NullString{String: util.RandomAddress(), Valid: true},
		CollegeID:   sql.NullInt64{Int64: college.CollegeID, Valid: true},
		FunnelID:    sql.NullInt64{Int64: funnel.FunnelID, Valid: true},
		HourlyFee:   sql.NullFloat64{Float64: util.RandomLessonHourlyFee(), Valid: true},
		Notes:       sql.NullString{String: util.RandomNote(), Valid: true},
	}

	student, err := testQueries.CreateStudent(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, student)

	require.Equal(t, arg.FirstName, student.FirstName)
	require.Equal(t, arg.LastName, student.LastName)
	require.Equal(t, arg.Email, student.Email)
	require.Equal(t, arg.PhoneNumber, student.PhoneNumber)
	require.Equal(t, arg.Address, student.Address)
	require.Equal(t, arg.CollegeID, student.CollegeID)
	require.Equal(t, arg.FunnelID, student.FunnelID)
	require.Equal(t, arg.HourlyFee, student.HourlyFee)
	require.Equal(t, arg.Notes, student.Notes)

	require.NotZero(t, student.StudentID)
	require.NotZero(t, student.CreatedAt)

	return student
}
func TestCreateStudent(t *testing.T) {
	createRandomStudent(t)
}

func TestGetStudent(t *testing.T) {
	student1 := createRandomStudent(t)
	student2, err := testQueries.GetStudent(context.Background(), student1.StudentID)
	require.NoError(t, err)
	require.NotEmpty(t, student2)

	require.Equal(t, student1.StudentID, student2.StudentID)
	require.Equal(t, student1.FirstName, student2.FirstName)
	require.Equal(t, student1.LastName, student2.LastName)
	require.Equal(t, student1.Email, student2.Email)
	require.Equal(t, student1.PhoneNumber, student2.PhoneNumber)
	require.Equal(t, student1.Address, student2.Address)
	require.Equal(t, student1.CollegeID, student2.CollegeID)
	require.Equal(t, student1.FunnelID, student2.FunnelID)
	require.Equal(t, student1.HourlyFee, student2.HourlyFee)
	require.Equal(t, student1.Notes, student2.Notes)
	require.WithinDuration(t, student1.CreatedAt, student2.CreatedAt, time.Second)
}

func TestUpdateStudent(t *testing.T) {
	student1 := createRandomStudent(t)
	college := createRandomCollege(t)
	funnel := createRandomFunnel(t)

	arg := UpdateStudentParams{
		StudentID:   student1.StudentID,
		FirstName:   util.RandomName(),
		LastName:    util.RandomName(),
		Email:       sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber: sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		Address:     sql.NullString{String: util.RandomAddress(), Valid: true},
		CollegeID:   sql.NullInt64{Int64: college.CollegeID, Valid: true},
		FunnelID:    sql.NullInt64{Int64: funnel.FunnelID, Valid: true},
		HourlyFee:   sql.NullFloat64{Float64: util.RandomLessonHourlyFee(), Valid: true},
		Notes:       sql.NullString{String: util.RandomNote(), Valid: true},
	}

	err := testQueries.UpdateStudent(context.Background(), arg)
	require.NoError(t, err)

	student2, err := testQueries.GetStudent(context.Background(), student1.StudentID)
	require.NoError(t, err)
	require.NotEmpty(t, student2)

	require.Equal(t, arg.StudentID, student2.StudentID)
	require.Equal(t, arg.FirstName, student2.FirstName)
	require.Equal(t, arg.LastName, student2.LastName)
	require.Equal(t, arg.Email, student2.Email)
	require.Equal(t, arg.PhoneNumber, student2.PhoneNumber)
	require.Equal(t, arg.Address, student2.Address)
	require.Equal(t, arg.CollegeID, student2.CollegeID)
	require.Equal(t, arg.FunnelID, student2.FunnelID)
	require.Equal(t, arg.HourlyFee, student2.HourlyFee)
	require.Equal(t, arg.Notes, student2.Notes)
	require.WithinDuration(t, student1.CreatedAt, student2.CreatedAt, time.Second)
}

func TestDeleteStudent(t *testing.T) {
	student1 := createRandomStudent(t)
	err := testQueries.DeleteStudent(context.Background(), student1.StudentID)
	require.NoError(t, err)

	student2, err := testQueries.GetStudent(context.Background(), student1.StudentID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, student2)
}

func TestListStudents(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomStudent(t)
	}

	arg := ListStudentsParams{
		Limit:  5,
		Offset: 5,
	}

	students, err := testQueries.ListStudents(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, students, 5)

	for _, student := range students {
		require.NotEmpty(t, student)
	}
}
