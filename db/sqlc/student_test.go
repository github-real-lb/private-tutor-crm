package db

import (
	"context"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

func TestCreateStudent(t *testing.T) {
	// CreateStudentParams with Null values
	arg := CreateStudentParams{
		FirstName:   util.RandomName(),
		LastName:    util.RandomName(),
		Email:       util.NullNullSting(),
		PhoneNumber: util.NullNullSting(),
		Address:     util.NullNullSting(),
		CollegeID:   util.NullNullInt64(),
		FunnelID:    util.NullNullInt64(),
		HourlyFee:   util.NullNullFloat64(),
		Notes:       util.NullNullSting(),
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

	// add record to college

	// add record to funnel

	// CreateStudentParams without Null values
}
