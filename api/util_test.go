package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/github-real-lb/tutor-management-web/db/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCase used as a single test case for APIs
type testCase struct {
	name          string      // name of test
	httpMethod    string      // http.Method for the http.Request
	url           string      // url for the http.Request
	body          interface{} // the json body for the http.Request
	buildStub     func(mockStore *mocks.MockStore)
	checkResponse func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder)
}

// sendRequestToTestServer start test server and send the test request
func (tc *testCase) sendRequestToServer(t *testing.T, mockStore *mocks.MockStore) *httptest.ResponseRecorder {
	// start test server and send request
	server := NewServer(mockStore)
	recorder := httptest.NewRecorder()

	var reader io.Reader = nil

	// creating new reader with arguments passed
	if tc.body != nil {
		jsonData, err := json.Marshal(tc.body)
		require.NoError(t, err)

		reader = strings.NewReader(string(jsonData))
	}

	request, err := http.NewRequest(tc.httpMethod, tc.url, reader)
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
