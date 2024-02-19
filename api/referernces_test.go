package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/github-real-lb/tutor-management-web/db/mock"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetReferenceStructsAPI(t *testing.T) {
	ref := randomReferenceStruct(db.ReferenceCollege)
	// TODO: loop through all references map and check all APIs once done.

	testCases := getReferenceStructTestCases(ref)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// start mock db and build the GetStudent stub
			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			// send test request to server
			url := fmt.Sprintf("/colleges/%d", tc.id)
			recorder := sendRequestToTestServer(t, store, http.MethodGet, url)

			tc.checkResponse(t, recorder)
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

// getCollegeTestCases generate a collection of tests for the getCollege API
func getReferenceStructTestCases(ref db.ReferenceStruct) []getTestCase {
	var testCases []getTestCase

	if ref == nil {
		return nil
	}

	refID := ref.GetID()
	refStructName := ref.GetReferenceStructName()

	// StatusOK API response test case
	testCases = append(testCases, getTestCase{
		name: "OK",
		id:   refID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetCollege(gomock.Any(), refID).
				Times(1).
				Return(ref, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, ref)
		},
	})

	// Record Not Found API response test case
	testCases = append(testCases, getTestCase{
		name: "NotFound",
		id:   refID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetCollege(gomock.Any(), refID).
				Times(1).
				Return(db.ReferenceStructMap[refStructName], sql.ErrNoRows)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusNotFound, recorder.Code)
		},
	})

	// Server Internal Error API response test case
	testCases = append(testCases, getTestCase{
		name: "InternalError",
		id:   refID,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetCollege(gomock.Any(), refID).
				Times(1).
				Return(db.ReferenceStructMap[refStructName], sql.ErrConnDone)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusInternalServerError, recorder.Code)
		},
	})

	// Invalid ID API response test case
	testCases = append(testCases, getTestCase{
		name: "InvalidID",
		id:   0,
		buildStub: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetCollege(gomock.Any(), gomock.Any()).
				Times(0)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusBadRequest, recorder.Code)
		},
	})

	return testCases
}
