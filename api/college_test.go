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

func TestCollegeAPIs(t *testing.T) {
	tests := tests{
		"Test_createCollegeAPI": createCollegeTestCasesBuilder(),
		"Test_getCollege":       getCollegeTestCasesBuilder(),
		"Test_listColleges":     listCollegesTestCasesBuilder(),
		"Test_updateCollege":    updateCollegeTestCasesBuilder(),
	}

	for key, tcs := range tests {
		t.Run(key, func(t *testing.T) {
			for _, tc := range tcs {
				t.Run(tc.name, func(t *testing.T) {
					// start mock db and build the stub
					mockStore := mocks.NewMockStore(t)
					tc.buildStub(mockStore)

					// send test request to server
					recorder := tc.sendRequestToServer(t, mockStore)

					// check response
					tc.checkResponse(t, mockStore, recorder)
				})
			}

		})
	}
}

// randomCollege creates a new random College struct.
func randomCollege() db.College {
	return db.College{
		CollegeID: util.RandomInt64(1, 1000),
		Name:      util.RandomName(),
	}
}

// createCollegeTestCasesBuilder creates a slice of test cases for the createCollege API
func createCollegeTestCasesBuilder() testCases {
	var testCases testCases

	college := randomCollege()

	arg := struct {
		Name string `json:"name"`
	}{
		Name: college.Name,
	}

	methodName := "CreateCollege"
	url := "/colleges"

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodPost,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, college.Name).
				Return(college, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, college)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodPost,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return(db.College{}, sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		},
	})

	// create a test case for Invalid Body Data response by passing no arguments
	testCases = append(testCases, testCase{
		name:       "Invalid Body Data",
		httpMethod: http.MethodPost,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()

		},
	})

	return testCases
}

// getCollegeTestCasesBuilder creates a slice of test cases for the getCollege API
func getCollegeTestCasesBuilder() testCases {
	var testCases testCases

	college := randomCollege()
	id := college.CollegeID
	methodName := "GetCollege"
	url := fmt.Sprintf("/colleges/%d", id)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(college, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, college)

		},
	})

	// create a test case for Not Found response
	testCases = append(testCases, testCase{
		name:       "Not Found",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(db.College{}, sql.ErrNoRows).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusNotFound, recorder.Code)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return(db.College{}, sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		},
	})

	// create a test case for Invalid ID response by passing url with id=0
	testCases = append(testCases, testCase{
		name:       "Invalid ID",
		httpMethod: http.MethodGet,
		url:        "/colleges/0",
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()

		},
	})

	return testCases
}

// listCollegesTestCasesBuilder creates a slice of test cases for the listColleges API
func listCollegesTestCasesBuilder() testCases {
	var testCases testCases

	n := 5
	colleges := make([]db.College, n)
	for i := 0; i < n; i++ {
		colleges[i] = randomCollege()
	}

	arg := db.ListCollegesParams{
		Limit:  int32(n),
		Offset: 0,
	}

	methodName := "ListColleges"
	url := fmt.Sprintf("/colleges?page_id=%d&page_size=%d", 1, n)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, arg).
				Return(colleges, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, colleges)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return([]db.College{}, sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		},
	})

	// create a test case for Invalid PageID response by passing url with page_id=-1
	testCases = append(testCases, testCase{
		name:       "Invalid Page_ID Parameter",
		httpMethod: http.MethodGet,
		url:        fmt.Sprintf("/colleges?page_id=%d&page_size=%d", -1, n),
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()

		},
	})

	// create a test case for Invalid PageSize response by passing url with page_size=10000
	testCases = append(testCases, testCase{
		name:       "Invalid Page_Size Parameter",
		httpMethod: http.MethodGet,
		url:        fmt.Sprintf("/colleges?page_id=%d&page_size=%d", 1, 10000),
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()

		},
	})

	return testCases
}

// updateCollegeTestCasesBuilder creates a slice of test cases for the updateCollege API
func updateCollegeTestCasesBuilder() testCases {
	var testCases testCases

	arg := db.UpdateCollegeParams{
		CollegeID: util.RandomInt64(1, 1000),
		Name:      util.RandomName(),
	}

	methodName := "UpdateCollege"
	url := "/colleges"

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodPut,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, arg).
				Return(nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodPut,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return(sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		},
	})

	// create a test case for Invalid Body Data response by passing no arguments
	testCases = append(testCases, testCase{
		name:       "Invalid Body Data",
		httpMethod: http.MethodPut,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()
		},
	})

	return testCases
}
