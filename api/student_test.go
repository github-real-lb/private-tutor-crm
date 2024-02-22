package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/github-real-lb/tutor-management-web/db/mocks"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetStudentAPI(t *testing.T) {
	testCases := testCasesGetStudentBuilder()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// start mock db and build the GetStudent stub
			mockStore := mocks.NewMockStore(t)
			tc.buildStub(mockStore)

			// send test request to server
			recorder := tc.sendRequestToServer(t, mockStore)

			// check response
			tc.checkResponse(t, mockStore, recorder)
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

func testCasesGetStudentBuilder() []testCaseGet {
	var testCases []testCaseGet

	student := randomStudent()
	id := student.StudentID
	methodName := "GetStudent"
	url := fmt.Sprintf("/students/%d", id)

	// create a test case for StatusOK response
	testCases = append(testCases, testCaseGet{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(student, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, student)
			mockStore.AssertExpectations(t)
		},
	})

	// create a test case for Not Found response
	testCases = append(testCases, testCaseGet{
		name:       "Not Found",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(&db.Student{}, sql.ErrNoRows).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusNotFound, recorder.Code)
			mockStore.AssertExpectations(t)
		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCaseGet{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(&db.Student{}, sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			mockStore.AssertExpectations(t)
		},
	})

	// create a test case for Invalid ID response by passing id=0
	testCases = append(testCases, testCaseGet{
		name:       "Invalid ID",
		httpMethod: http.MethodGet,
		url:        "/students/0",
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()
			mockStore.AssertNotCalled(t, methodName)
		},
	})

	return testCases
}
