package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/stretchr/testify/require"
)

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
