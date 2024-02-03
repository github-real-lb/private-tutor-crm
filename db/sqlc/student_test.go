package db

import (
	"context"
	"database/sql"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const ABC = "abcdefghijklmnopqrstuvwxyz"

func TestCreateStudent(t *testing.T) {
	// CreateStudentParams with Null values
	arg := CreateStudentParams{
		FirstName:   GetRandomString(10),
		LastName:    GetRandomString(10),
		Email:       GetRandomNullString(20, true),
		PhoneNumber: GetRandomNullString(10, true),
		Address:     GetRandomNullString(10, true),
		CollegeID:   GetRandomNullInt64(0, false),
		FunnelID:    GetRandomNullInt64(0, false),
		HourlyFee:   GetRandomNullFloat64(0, false),
		Notes:       GetRandomNullString(50, true),
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
	collegeID := sql.NullInt64{
		Int64: 1,
		Valid: true,
	}
	// add record to funnel
	funnelID := sql.NullInt64{
		Int64: 1,
		Valid: true,
	}

	// CreateStudentParams without Null values
	arg = CreateStudentParams{
		FirstName:   GetRandomString(10),
		LastName:    GetRandomString(10),
		Email:       GetRandomNullString(20, true),
		PhoneNumber: GetRandomNullString(10, true),
		Address:     GetRandomNullString(10, true),
		CollegeID:   collegeID,
		FunnelID:    funnelID,
		HourlyFee:   GetRandomNullFloat64(5, true),
		Notes:       GetRandomNullString(50, true),
	}
}

func GetRandomString(n int) string {
	var sb strings.Builder
	var b byte

	for i := 0; i < n-1; i++ {
		b = ABC[rand.Intn(len(ABC))]
		sb.WriteByte(b)
	}

	return sb.String()
}

func GetRandomInt64(n int) int64 {
	return int64(rand.Intn(n))
}

func GetRandomFloat64(n int) float64 {
	return float64(rand.Intn(n)) / 100
}

func GetRandomNullString(n int, valid bool) sql.NullString {
	if valid {
		return sql.NullString{
			String: GetRandomString(n),
			Valid:  true,
		}
	} else {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
}

func GetRandomNullInt64(n int, valid bool) sql.NullInt64 {
	if valid {
		return sql.NullInt64{
			Int64: GetRandomInt64(n),
			Valid: true,
		}
	} else {
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}
}

func GetRandomNullFloat64(n int, valid bool) sql.NullFloat64 {
	if valid {
		return sql.NullFloat64{
			Float64: GetRandomFloat64(n),
			Valid:   true,
		}
	} else {
		return sql.NullFloat64{
			Float64: 0.0,
			Valid:   false,
		}
	}
}
