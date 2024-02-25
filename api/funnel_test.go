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

func TestFunnelAPIs(t *testing.T) {
	tests := tests{
		"Test_createFunnelAPI": createFunnelTestCasesBuilder(),
		"Test_getFunnel":       getFunnelTestCasesBuilder(),
		"Test_listFunnels":     listFunnelsTestCasesBuilder(),
		"Test_updateFunnel":    updateFunnelTestCasesBuilder(),
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

// randomFunnel creates a new random Funnel struct.
func randomFunnel() db.Funnel {
	return db.Funnel{
		FunnelID: util.RandomInt64(1, 1000),
		Name:     util.RandomName(),
	}
}

// createFunnelTestCasesBuilder creates a slice of test cases for the createFunnel API
func createFunnelTestCasesBuilder() testCases {
	var testCases testCases

	funnel := randomFunnel()

	arg := struct {
		Name string `json:"name"`
	}{
		Name: funnel.Name,
	}

	methodName := "CreateFunnel"
	url := "/funnels"

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodPost,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, funnel.Name).
				Return(funnel, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, funnel)

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
				Return(db.Funnel{}, sql.ErrConnDone).
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

// getFunnelTestCasesBuilder creates a slice of test cases for the getFunnel API
func getFunnelTestCasesBuilder() testCases {
	var testCases testCases

	funnel := randomFunnel()
	id := funnel.FunnelID
	methodName := "GetFunnel"
	url := fmt.Sprintf("/funnels/%d", id)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(funnel, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, funnel)

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
				Return(db.Funnel{}, sql.ErrNoRows).
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
				Return(db.Funnel{}, sql.ErrConnDone).
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
		url:        "/funnels/0",
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

// listFunnelsTestCasesBuilder creates a slice of test cases for the listFunnels API
func listFunnelsTestCasesBuilder() testCases {
	var testCases testCases

	n := 5
	funnels := make([]db.Funnel, n)
	for i := 0; i < n; i++ {
		funnels[i] = randomFunnel()
	}

	arg := db.ListFunnelsParams{
		Limit:  int32(n),
		Offset: 0,
	}

	methodName := "ListFunnels"
	url := fmt.Sprintf("/funnels?page_id=%d&page_size=%d", 1, n)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, arg).
				Return(funnels, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, funnels)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return([]db.Funnel{}, sql.ErrConnDone).
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
		url:        fmt.Sprintf("/funnels?page_id=%d&page_size=%d", -1, n),
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
		url:        fmt.Sprintf("/funnels?page_id=%d&page_size=%d", 1, 10000),
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

// updateFunnelTestCasesBuilder creates a slice of test cases for the updateFunnel API
func updateFunnelTestCasesBuilder() testCases {
	var testCases testCases

	arg := db.UpdateFunnelParams{
		FunnelID: util.RandomInt64(1, 1000),
		Name:     util.RandomName(),
	}

	methodName := "UpdateFunnel"
	url := "/funnels"

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
