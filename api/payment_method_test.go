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

func TestPaymentMethodAPIs(t *testing.T) {
	tests := tests{
		"Test_createPaymentMethodAPI": createPaymentMethodTestCasesBuilder(),
		"Test_getPaymentMethod":       getPaymentMethodTestCasesBuilder(),
		"Test_listPaymentMethods":     listPaymentMethodsTestCasesBuilder(),
		"Test_updatePaymentMethods":   updatePaymentMethodTestCasesBuilder(),
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

// randomPaymentMethod creates a new random PaymentMethod struct.
func randomPaymentMethod() db.PaymentMethod {
	return db.PaymentMethod{
		PaymentMethodID: util.RandomInt64(1, 1000),
		Name:            util.RandomName(),
	}
}

// createPaymentMethodTestCasesBuilder creates a slice of test cases for the createPaymentMethod API
func createPaymentMethodTestCasesBuilder() testCases {
	var testCases testCases

	paymentMethod := randomPaymentMethod()

	arg := struct {
		Name string `json:"name"`
	}{
		Name: paymentMethod.Name,
	}

	methodName := "CreatePaymentMethod"
	url := "/payment_methods"

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodPost,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, paymentMethod.Name).
				Return(paymentMethod, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, paymentMethod)

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
				Return(db.PaymentMethod{}, sql.ErrConnDone).
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

// getPaymentMethodTestCasesBuilder creates a slice of test cases for the getPaymentMethod API
func getPaymentMethodTestCasesBuilder() testCases {
	var testCases testCases

	paymentMethod := randomPaymentMethod()
	id := paymentMethod.PaymentMethodID
	methodName := "GetPaymentMethod"
	url := fmt.Sprintf("/payment_methods/%d", id)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(paymentMethod, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, paymentMethod)

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
				Return(db.PaymentMethod{}, sql.ErrNoRows).
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
				Return(db.PaymentMethod{}, sql.ErrConnDone).
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
		url:        "/payment_methods/0",
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

// listPaymentMethodsTestCasesBuilder creates a slice of test cases for the listPaymentMethods API
func listPaymentMethodsTestCasesBuilder() testCases {
	var testCases testCases

	n := 5
	paymentMethods := make([]db.PaymentMethod, n)
	for i := 0; i < n; i++ {
		paymentMethods[i] = randomPaymentMethod()
	}

	arg := db.ListPaymentMethodsParams{
		Limit:  int32(n),
		Offset: 0,
	}

	methodName := "ListPaymentMethods"
	url := fmt.Sprintf("/payment_methods?page_id=%d&page_size=%d", 1, n)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, arg).
				Return(paymentMethods, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, paymentMethods)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return([]db.PaymentMethod{}, sql.ErrConnDone).
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
		url:        fmt.Sprintf("/payment_methods?page_id=%d&page_size=%d", -1, n),
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
		url:        fmt.Sprintf("/payment_methods?page_id=%d&page_size=%d", 1, 10000),
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

// updatePaymentMethodTestCasesBuilder creates a slice of test cases for the updatePaymentMethod API
func updatePaymentMethodTestCasesBuilder() testCases {
	var testCases testCases

	arg := db.UpdatePaymentMethodParams{
		PaymentMethodID: util.RandomInt64(1, 1000),
		Name:            util.RandomName(),
	}

	methodName := "UpdatePaymentMethod"
	url := "/payment_methods"

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
