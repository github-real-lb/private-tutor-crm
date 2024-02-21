package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/github-real-lb/tutor-management-web/db/mocks"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAPI(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T, name db.ReferenceStructName)
	}{
		{"OK", testGetAPI_OK},
		{"Not Found", testGetAPI_NotFound},
		{"Internal Error", testGetAPI_InternalError},
		{"Invalid ID", testGetAPI_InvalidID},
	}

	for key := range db.ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					tc.test(t, key)
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

func testGetAPI_OK(t *testing.T, name db.ReferenceStructName) {
	ref := randomReferenceStruct(name)

	// create mockStore
	mockStore := mocks.NewMockStore(t)

	// create stub
	methodName := fmt.Sprintf("Get%s", string(name))
	mockStore.On(methodName, mock.Anything, ref.GetID()).
		Return(ref, nil).
		Once()

	// send test request to server
	url := fmt.Sprintf("/%ss/%d", strings.ToLower(string(name)), ref.GetID())
	recorder := sendRequestToTestServer(t, mockStore, http.MethodGet, url)

	// check server response
	assert.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchStruct(t, recorder.Body, ref)

	mockStore.AssertExpectations(t)
}

func testGetAPI_NotFound(t *testing.T, name db.ReferenceStructName) {
	ref := randomReferenceStruct(name)

	// create mockStore
	mockStore := mocks.NewMockStore(t)

	// create stub
	methodName := fmt.Sprintf("Get%s", string(name))
	mockStore.On(methodName, mock.Anything, ref.GetID()).
		Return(db.ReferenceStructMap[ref.GetReferenceStructName()], sql.ErrNoRows).
		Once()

	// send test request to server
	url := fmt.Sprintf("/%ss/%d", strings.ToLower(string(name)), ref.GetID())
	recorder := sendRequestToTestServer(t, mockStore, http.MethodGet, url)

	// check server response
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	mockStore.AssertExpectations(t)
}

func testGetAPI_InternalError(t *testing.T, name db.ReferenceStructName) {
	ref := randomReferenceStruct(name)

	// create mockStore
	mockStore := mocks.NewMockStore(t)

	// create stub
	methodName := fmt.Sprintf("Get%s", string(name))
	mockStore.On(methodName, mock.Anything, ref.GetID()).
		Return(db.ReferenceStructMap[ref.GetReferenceStructName()], sql.ErrConnDone).
		Once()

	// send test request to server
	url := fmt.Sprintf("/%ss/%d", strings.ToLower(string(name)), ref.GetID())
	recorder := sendRequestToTestServer(t, mockStore, http.MethodGet, url)

	// check server response
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	mockStore.AssertExpectations(t)
}

func testGetAPI_InvalidID(t *testing.T, name db.ReferenceStructName) {
	// create mockStore
	mockStore := mocks.NewMockStore(t)

	// create stub
	methodName := fmt.Sprintf("Get%s", string(name))
	mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)

	// send test request to server
	url := fmt.Sprintf("/%ss/%d", strings.ToLower(string(name)), 0)
	recorder := sendRequestToTestServer(t, mockStore, http.MethodGet, url)

	// check server response
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	mockStore.On(methodName, mock.Anything, mock.Anything).Unset()
	mockStore.AssertNotCalled(t, methodName)
}
