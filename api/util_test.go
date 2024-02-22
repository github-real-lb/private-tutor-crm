package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/github-real-lb/tutor-management-web/db/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTestCase used as a single test case for the get record API
type testCaseGet struct {
	name          string
	httpMethod    string
	url           string
	buildStub     func(mockStore *mocks.MockStore)
	checkResponse func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder)
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
