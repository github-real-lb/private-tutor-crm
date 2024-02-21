package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/github-real-lb/tutor-management-web/db/mocks"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// getTestCase used as a single test case for the get record API
type testCaseGet struct {
	name          string
	httpMethod    string
	url           string
	id            int64 // record id to test
	buildStub     func(mockStore *mocks.MockStore)
	checkResponse func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder)
}

func TestGetAPI(t *testing.T) {
	for key := range db.ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			testCases := testCasesGetBuilder(key)
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

func randomReferenceStruct(key db.ReferenceStructName) db.ReferenceStruct {
	ref, ok := db.ReferenceStructMap[key]
	if !ok {
		return nil
	}

	ref.SetID(util.RandomInt64(1, 1000))
	ref.SetName(util.RandomName())

	return ref
}

func testCasesGetBuilder(name db.ReferenceStructName) []testCaseGet {
	var testCases []testCaseGet

	ref := randomReferenceStruct(name)
	id := ref.GetID()

	methodName := fmt.Sprintf("Get%s", string(name))
	url := fmt.Sprintf("/%ss/%d", strings.ToLower(string(name)), id)

	// create a test case for StatusOK response
	testCases = append(testCases, testCaseGet{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		id:         id,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, ref.GetID()).
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
	testCases = append(testCases, testCaseGet{
		name:       "Not Found",
		httpMethod: http.MethodGet,
		url:        url,
		id:         id,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, ref.GetID()).
				Return(db.ReferenceStructMap[ref.GetReferenceStructName()], sql.ErrNoRows).
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
		id:         id,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, ref.GetID()).
				Return(db.ReferenceStructMap[ref.GetReferenceStructName()], sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			mockStore.AssertExpectations(t)
		},
	})

	// create a test case for Invalid ID response
	testCases = append(testCases, testCaseGet{
		name:       "Invalid ID",
		httpMethod: http.MethodGet,
		url:        fmt.Sprintf("/%ss/%d", strings.ToLower(string(name)), 0),
		id:         0,
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

// sendRequestToTestServer start test server and send the test request
func (tc *testCaseGet) sendRequestToServer(t *testing.T, mockStore *mocks.MockStore) *httptest.ResponseRecorder {
	// start test server and send request
	server := NewServer(mockStore)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(tc.httpMethod, tc.url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	return recorder
}

// requireBodyMatchStruct asserts that a JSON httptest.ResponseRecorder.Body
// equal to a Struct object.
func requireBodyMatchStruct(t *testing.T, body *bytes.Buffer, obj interface{}) {
	jsonBodyData, err := io.ReadAll(body)
	require.NoError(t, err)

	jsonObjData, err := json.Marshal(obj)
	require.NoError(t, err)
	assert.Equal(t, string(jsonObjData), string(jsonBodyData))
}
