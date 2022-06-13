package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Pythonyan3/payment-service/internal/models"
	mock_services "github.com/Pythonyan3/payment-service/internal/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_TransactionsListByUserId(t *testing.T) {
	// Arrange
	type mockBehaviour func(service *mock_services.MockUserService, userId int)
	serializedTransactions, _ := json.Marshal(transactionSlice)

	testTable := []struct {
		name                string
		userId              int
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Test list of transactions (ok)",
			userId:              transaction.UserId,
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(serializedTransactions) + "\n",
			mockBehaviour: func(service *mock_services.MockUserService, userId int) {
				service.EXPECT().GetUserTransactionsById(userId).Return(transactionSlice, nil)
			},
		},
		{
			name:                "Test list of transactions (empty)",
			userId:              transaction.UserId,
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "[]\n",
			mockBehaviour: func(service *mock_services.MockUserService, userId int) {
				service.EXPECT().GetUserTransactionsById(userId).Return([]*models.Transaction{}, nil)
			},
		},
		{
			name:                "Test list of transactions (service error)",
			userId:              transaction.UserId,
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: "Internal Server Error\n",
			mockBehaviour: func(service *mock_services.MockUserService, userId int) {
				service.EXPECT().GetUserTransactionsById(userId).Return(nil, errors.New("some error"))
			},
		},
	}

	// Act
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_services.NewMockUserService(controller)
			testCase.mockBehaviour(service, testCase.userId)

			handler := NewUserHandler(service)
			router := mux.NewRouter()
			router.HandleFunc("/api/users/{userId:[0-9]+}/transactions/", handler.TransactionsListByUserId)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				"GET", fmt.Sprintf("/api/users/%d/transactions/", testCase.userId), bytes.NewBufferString(""),
			)

			router.ServeHTTP(w, r)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_TransactionsListByUserEmail(t *testing.T) {
	// Arrange
	type mockBehaviour func(service *mock_services.MockUserService, userEmail string)
	serializedTransactions, _ := json.Marshal(transactionSlice)

	testTable := []struct {
		name                string
		userEmail           string
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Test list of transactions (ok)",
			userEmail:           transaction.UserEmail,
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(serializedTransactions) + "\n",
			mockBehaviour: func(service *mock_services.MockUserService, userEmail string) {
				service.EXPECT().GetUserTransactionsByEmail(userEmail).Return(transactionSlice, nil)
			},
		},
		{
			name:                "Test list of transactions (empty)",
			userEmail:           transaction.UserEmail,
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "[]\n",
			mockBehaviour: func(service *mock_services.MockUserService, userEmail string) {
				service.EXPECT().GetUserTransactionsByEmail(userEmail).Return([]*models.Transaction{}, nil)
			},
		},
		{
			name:                "Test list of transactions (service error)",
			userEmail:           transaction.UserEmail,
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: "Internal Server Error\n",
			mockBehaviour: func(service *mock_services.MockUserService, userEmail string) {
				service.EXPECT().GetUserTransactionsByEmail(userEmail).Return(nil, errors.New("some error"))
			},
		},
	}

	// Act
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_services.NewMockUserService(controller)
			testCase.mockBehaviour(service, testCase.userEmail)

			handler := NewUserHandler(service)
			router := mux.NewRouter()
			router.HandleFunc("/api/users/{userEmail}/transactions/", handler.TransactionsListByUserEmail)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				"GET", fmt.Sprintf("/api/users/%s/transactions/", testCase.userEmail), bytes.NewBufferString(""),
			)

			router.ServeHTTP(w, r)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
