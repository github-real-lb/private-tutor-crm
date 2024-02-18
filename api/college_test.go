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

func TestGetCollegeAPI(t *testing.T) {
	college := randomCollege()

	testCases := getCollegeTestCases(college)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// start mock db and build the GetStudent stub
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			// send test request to server
			url := fmt.Sprintf("/colleges/%d", tc.id)
			recorder := sendRequestToTestServer(t, store, http.MethodGet, url)

			tc.checkResponse(t, recorder)
		})
	}

}

func randomCollege() db.College {
	return db.College{
		CollegeID: util.RandomInt64(1, 1000),
		Name:      util.RandomName(),
	}
}

// getCollegeTestCases generate a collection of tests for the getCollege API
func getCollegeTestCases(college db.College) []getTestCase {
	var testCases []getTestCase

	// StatusOK API response test case
	testCases = append(testCases, getTestCase{
		name: "OK",
		id:   college.CollegeID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetCollege(gomock.Any(), college.CollegeID).
				Times(1).
				Return(college, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, college)
		},
	})

	// Record Not Found API response test case
	testCases = append(testCases, getTestCase{
		name: "NotFound",
		id:   college.CollegeID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetCollege(gomock.Any(), college.CollegeID).
				Times(1).
				Return(db.College{}, sql.ErrNoRows)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusNotFound, recorder.Code)
		},
	})

	// Server Internal Error API response test case
	testCases = append(testCases, getTestCase{
		name: "InternalError",
		id:   college.CollegeID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetCollege(gomock.Any(), college.CollegeID).
				Times(1).
				Return(db.College{}, sql.ErrConnDone)
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
				GetCollege(gomock.Any(), gomock.Any()).
				Times(0)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusBadRequest, recorder.Code)
		},
	})

	return testCases
}
