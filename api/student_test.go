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

	testCases := getStudentTestCases(student)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// start mock db and build the GetStudent stub
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			// send test request to server
			url := fmt.Sprintf("/students/%d", tc.id)
			recorder := sendRequestToTestServer(t, store, http.MethodGet, url)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomStudent() *db.Student {
	return &db.Student{
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

// getStudentTestCases generate a collection of tests for the getStudent API
func getStudentTestCases(student *db.Student) []getTestCase {
	var testCases []getTestCase

	// StatusOK API response test case
	testCases = append(testCases, getTestCase{
		name: "OK",
		id:   student.StudentID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetStudent(gomock.Any(), student.StudentID).
				Times(1).
				Return(student, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, student)
		},
	})

	// Record Not Found API response test case
	testCases = append(testCases, getTestCase{
		name: "NotFound",
		id:   student.StudentID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetStudent(gomock.Any(), student.StudentID).
				Times(1).
				Return(&db.Student{}, sql.ErrNoRows)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusNotFound, recorder.Code)
		},
	})

	// Server Internal Error API response test case
	testCases = append(testCases, getTestCase{
		name: "InternalError",
		id:   student.StudentID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetStudent(gomock.Any(), student.StudentID).
				Times(1).
				Return(&db.Student{}, sql.ErrConnDone)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusInternalServerError, recorder.Code)
		},
	})

	// Invalid ID API response test case
	testCases = append(testCases, getTestCase{
		name: "InvalidID",
		id:   0,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetStudent(gomock.Any(), gomock.Any()).
				Times(0)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusBadRequest, recorder.Code)
		},
	})

	return testCases
}

// testCases := []struct {
// 	name          string
// 	studentID     int64
// 	buildStubs    func(store *mockdb.MockStore)
// 	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// }{
// 	{
// 		name:      "OK",
// 		studentID: student.StudentID,
// 		buildStubs: func(store *mockdb.MockStore) {
// 			store.EXPECT().
// 				GetStudent(gomock.Any(), student.StudentID).
// 				Times(1).
// 				Return(student, nil)
// 		},
// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			require.Equal(t, http.StatusOK, recorder.Code)
// 			requireBodyMatchStruct(t, recorder.Body, student)
// 		},
// 	},
// 	{
// 		name:      "NotFound",
// 		studentID: student.StudentID,
// 		buildStubs: func(store *mockdb.MockStore) {
// 			store.EXPECT().
// 				GetStudent(gomock.Any(), student.StudentID).
// 				Times(1).
// 				Return(db.Student{}, sql.ErrNoRows)
// 		},
// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			require.Equal(t, http.StatusNotFound, recorder.Code)
// 		},
// 	},
// 	{
// 		name:      "InternalError",
// 		studentID: student.StudentID,
// 		buildStubs: func(store *mockdb.MockStore) {
// 			store.EXPECT().
// 				GetStudent(gomock.Any(), student.StudentID).
// 				Times(1).
// 				Return(db.Student{}, sql.ErrConnDone)
// 		},
// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 		},
// 	},
// 	{
// 		name:      "InvalidID",
// 		studentID: 0,
// 		buildStubs: func(store *mockdb.MockStore) {
// 			store.EXPECT().
// 				GetStudent(gomock.Any(), gomock.Any()).
// 				Times(0)
// 		},
// 		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			require.Equal(t, http.StatusBadRequest, recorder.Code)
// 		},
// 	},
// }
