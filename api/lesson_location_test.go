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

func TestLessonLocationAPIs(t *testing.T) {
	tests := tests{
		"Test_createLessonLocationAPI": createLessonLocationTestCasesBuilder(),
		"Test_getLessonLocation":       getLessonLocationTestCasesBuilder(),
		"Test_listLessonLocations":     listLessonLocationsTestCasesBuilder(),
		"Test_updateLessonLocations":   updateLessonLocationTestCasesBuilder(),
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

// randomLessonLocation creates a new random LessonLocation struct.
func randomLessonLocation() db.LessonLocation {
	return db.LessonLocation{
		LocationID: util.RandomInt64(1, 1000),
		Name:       util.RandomName(),
	}
}

// createLessonLocationTestCasesBuilder creates a slice of test cases for the createLessonLocation API
func createLessonLocationTestCasesBuilder() testCases {
	var testCases testCases

	lessonLocation := randomLessonLocation()

	arg := struct {
		Name string `json:"name"`
	}{
		Name: lessonLocation.Name,
	}

	methodName := "CreateLessonLocation"
	url := "/lesson_locations"

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodPost,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, lessonLocation.Name).
				Return(lessonLocation, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, lessonLocation)

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
				Return(db.LessonLocation{}, sql.ErrConnDone).
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

// getLessonLocationTestCasesBuilder creates a slice of test cases for the getLessonLocation API
func getLessonLocationTestCasesBuilder() testCases {
	var testCases testCases

	lessonLocation := randomLessonLocation()
	id := lessonLocation.LocationID
	methodName := "GetLessonLocation"
	url := fmt.Sprintf("/lesson_locations/%d", id)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(lessonLocation, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, lessonLocation)

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
				Return(db.LessonLocation{}, sql.ErrNoRows).
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
				Return(db.LessonLocation{}, sql.ErrConnDone).
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
		url:        "/lesson_locations/0",
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

// listLessonLocationsTestCasesBuilder creates a slice of test cases for the listLessonLocations API
func listLessonLocationsTestCasesBuilder() testCases {
	var testCases testCases

	n := 5
	lessonLocations := make([]db.LessonLocation, n)
	for i := 0; i < n; i++ {
		lessonLocations[i] = randomLessonLocation()
	}

	arg := db.ListLessonLocationsParams{
		Limit:  int32(n),
		Offset: 0,
	}

	methodName := "ListLessonLocations"
	url := fmt.Sprintf("/lesson_locations?page_id=%d&page_size=%d", 1, n)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, arg).
				Return(lessonLocations, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, lessonLocations)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return([]db.LessonLocation{}, sql.ErrConnDone).
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
		url:        fmt.Sprintf("/lesson_locations?page_id=%d&page_size=%d", -1, n),
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
		url:        fmt.Sprintf("/lesson_locations?page_id=%d&page_size=%d", 1, 10000),
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

// updateLessonLocationTestCasesBuilder creates a slice of test cases for the updateLessonLocation API
func updateLessonLocationTestCasesBuilder() testCases {
	var testCases testCases

	arg := db.UpdateLessonLocationParams{
		LocationID: util.RandomInt64(1, 1000),
		Name:       util.RandomName(),
	}

	methodName := "UpdateLessonLocation"
	url := "/lesson_locations"

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
