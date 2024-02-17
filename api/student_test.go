package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/github-real-lb/tutor-management-web/db/mock"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetStudentAPI(t *testing.T) {
	student := randomStudent()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	//build stubs
	store.EXPECT().GetStudent(gomock.Any(), student.StudentID).Times(1).Return(student, nil)

	// start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/students/%d", student.StudentID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	// check the response
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchStruct(t, recorder.Body, student)

}

func randomStudent() db.Student {
	return db.Student{
		StudentID:   util.RandomInt64(1, 1000),
		FirstName:   util.RandomName(),
		LastName:    util.RandomName(),
		Email:       sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber: sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		Address:     sql.NullString{String: util.RandomAddress(), Valid: true},
		CollegeID:   sql.NullInt64{Int64: 0, Valid: false},
		FunnelID:    sql.NullInt64{Int64: 0, Valid: false},
		HourlyFee:   sql.NullFloat64{Float64: util.RandomLessonHourlyFee(), Valid: true},
		Notes:       sql.NullString{String: util.RandomNote(), Valid: true},
	}
}
