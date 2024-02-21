package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/github-real-lb/tutor-management-web/db/mock"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCaseType string

const (
	statusOK      testCaseType = "Test OK"
	notFound      testCaseType = "Test Not Found"
	internalError testCaseType = "Test Internal Server Error"
	invalidID     testCaseType = "Test Invalid ID"
)

// getTestCase used as a single test case for the get record API
type testCaseGet struct {
	name          testCaseType
	id            int64 // record id to test
	buildStub     func(store *mockdb.MockStore)
	sendRequest   func(t *testing.T, store db.Store) *httptest.ResponseRecorder
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

// sendRequestToTestServer start test server and send the test request
func sendRequestToTestServer(t *testing.T, store db.Store, httpMethod string, url string) *httptest.ResponseRecorder {
	// start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, url, nil)
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
