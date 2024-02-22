package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/github-real-lb/tutor-management-web/db/mocks"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// randomReferenceStruct creates a new random ReferenceStruct
// (College, Funnel, LessonLocation, LessonSubject or PaymentMethod).
func randomReferenceStruct(key db.ReferenceStructName) db.ReferenceStruct {
	ref, ok := db.ReferenceStructMap[key]
	if !ok {
		return nil
	}

	ref.SetID(util.RandomInt64(1, 1000))
	ref.SetName(util.RandomName())

	return ref
}

// createReferenceTestCasesBuilder creates a slice of test cases for the createReference API
// (createCollege, createFunnel, createLessonLocation, createLessonSubject or createPaymentMethod).
func createReferenceTestCasesBuilder(name db.ReferenceStructName) []testCase {
	var testCases []testCase

	ref := randomReferenceStruct(name)

	methodName := fmt.Sprintf("Create%s", string(name))
	url := fmt.Sprintf("/%ss", strings.ToLower(string(name)))

	arg := struct {
		Name string `json:"Name"`
	}{
		Name: ref.GetName(),
	}

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodPost,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, ref.GetName()).
				Return(ref, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, ref)
			mockStore.AssertExpectations(t)
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
				Return(db.ReferenceStructMap[name], sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			mockStore.AssertExpectations(t)
		},
	})

	// create a test case for Invalid arguments response by passing nil body
	testCases = append(testCases, testCase{
		name:       "Invalid Arguments",
		httpMethod: http.MethodPost,
		url:        url,
		body:       nil,
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

func Test_createReferenceAPI(t *testing.T) {
	for key := range db.ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			testCases := createReferenceTestCasesBuilder(key)
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
		})
	}
}

// getReferenceTestCasesBuilder creates a slice of test cases for the getReference API
// (getCollege, getFunnel, getLessonLocation, getLessonSubject or getPaymentMethod).
func getReferenceTestCasesBuilder(name db.ReferenceStructName) []testCase {
	var testCases []testCase

	ref := randomReferenceStruct(name)
	id := ref.GetID()

	methodName := fmt.Sprintf("Get%s", string(name))
	url := fmt.Sprintf("/%ss/%d", strings.ToLower(string(name)), id)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(ref, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, ref)
			mockStore.AssertExpectations(t)
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
				Return(db.ReferenceStructMap[name], sql.ErrNoRows).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusNotFound, recorder.Code)
			mockStore.AssertExpectations(t)
		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(db.ReferenceStructMap[name], sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			mockStore.AssertExpectations(t)
		},
	})

	// create a test case for Invalid ID response by passing id=0
	testCases = append(testCases, testCase{
		name:       "Invalid ID",
		httpMethod: http.MethodGet,
		url:        fmt.Sprintf("/%ss/%d", strings.ToLower(string(name)), 0),
		body:       nil,
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

func Test_getReferenceAPI(t *testing.T) {
	for key := range db.ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			testCases := getReferenceTestCasesBuilder(key)
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
		})
	}
}
