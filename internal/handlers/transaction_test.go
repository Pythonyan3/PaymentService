package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pythonyan3/payment-service/internal/models"
	"github.com/Pythonyan3/payment-service/internal/services"
	mock_services "github.com/Pythonyan3/payment-service/internal/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	currentTime time.Time           = time.Now()
	transaction *models.Transaction = &models.Transaction{
		Id:        1,
		UserId:    1,
		UserEmail: "email@mail.ru",
		Amount:    1500,
		Currency:  "EUR",
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Status:    services.TransactionNewStatus,
	}
	canceledTransaction *models.Transaction = &models.Transaction{
		Id:        transaction.Id,
		UserId:    transaction.UserId,
		UserEmail: transaction.UserEmail,
		Amount:    transaction.Amount,
		Currency:  transaction.Currency,
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
		Status:    transaction.Status,
	}
	transactionSlice []*models.Transaction    = []*models.Transaction{transaction}
	inputTransaction *models.TransactionInput = &models.TransactionInput{
		UserId:    transaction.Id,
		UserEmail: transaction.UserEmail,
		Amount:    transaction.Amount,
		Currency:  transaction.Currency,
	}
	badInputTransaction *models.TransactionInput = &models.TransactionInput{
		UserId:    transaction.Id,
		UserEmail: "not email",
		Amount:    transaction.Amount,
		Currency:  transaction.Currency,
	}
	emptyBody []byte = []byte{}
)

func TestHandler_RetrieveTransaction(t *testing.T) {
	// Arrange
	type mockBehaviour func(service *mock_services.MockTransactionService, transactionId int)
	serializedTransactions, _ := json.Marshal(transaction)

	testTable := []struct {
		name                string
		transactionId       int
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Test retrieve transaction (ok)",
			transactionId:       transaction.Id,
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(serializedTransactions) + "\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, transactionId int) {
				service.EXPECT().GetById(transactionId).Return(transaction, nil)
			},
		},
		{
			name:                "Test retrieve transaction (not found)",
			transactionId:       transaction.Id,
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: "Not Found\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, transactionId int) {
				service.EXPECT().GetById(transactionId).Return(nil, errors.New(dbNotFoundErrorMsg))
			},
		},
		{
			name:                "Test retrieve transaction (service error)",
			transactionId:       transaction.Id,
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: "Internal Server Error\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, transactionId int) {
				service.EXPECT().GetById(transactionId).Return(nil, errors.New("some error"))
			},
		},
	}

	// Act
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_services.NewMockTransactionService(controller)
			auth_service := mock_services.NewMockAuthMiddleware(controller)
			testCase.mockBehaviour(service, testCase.transactionId)

			handler := NewTransactionHandler(service, auth_service)
			router := mux.NewRouter()
			router.HandleFunc("/api/transactions/{pk:[0-9]+}/", handler.RetrieveTransaction)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				"GET", fmt.Sprintf("/api/transactions/%d/", testCase.transactionId), bytes.NewBufferString(""),
			)

			router.ServeHTTP(w, r)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_CreateTransaction(t *testing.T) {
	// Arrange
	type mockBehaviour func(service *mock_services.MockTransactionService, inputTransaction *models.TransactionInput)
	serializedTransaction, _ := json.Marshal(transaction)
	serializedInputTransaction, _ := json.Marshal(inputTransaction)
	serializedBadInputTransaction, _ := json.Marshal(badInputTransaction)

	testTable := []struct {
		name                string
		inputTransaction    *models.TransactionInput
		requestBody         []byte
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Test create transaction (ok)",
			inputTransaction:    inputTransaction,
			expectedStatusCode:  http.StatusCreated,
			requestBody:         serializedInputTransaction,
			expectedRequestBody: string(serializedTransaction) + "\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, inputTransaction *models.TransactionInput) {
				service.EXPECT().Create(inputTransaction).Return(transaction, nil)
			},
		},
		{
			name:                "Test create transaction (service error)",
			inputTransaction:    inputTransaction,
			requestBody:         serializedInputTransaction,
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: "Internal Server Error\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, inputTransaction *models.TransactionInput) {
				service.EXPECT().Create(inputTransaction).Return(nil, errors.New("some error"))
			},
		},
		{
			name:                "Test create transaction (no body)",
			inputTransaction:    inputTransaction,
			requestBody:         emptyBody,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "invalid request: EOF\n",
			mockBehaviour:       func(service *mock_services.MockTransactionService, inputTransaction *models.TransactionInput) {},
		},
		{
			name:                "Test create transaction (bad body)",
			inputTransaction:    badInputTransaction,
			requestBody:         serializedBadInputTransaction,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "invalid request: Key: 'TransactionInput.UserEmail' Error:Field validation for 'UserEmail' failed on the 'email' tag\n",
			mockBehaviour:       func(service *mock_services.MockTransactionService, inputTransaction *models.TransactionInput) {},
		},
	}

	// Act
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_services.NewMockTransactionService(controller)
			auth_service := mock_services.NewMockAuthMiddleware(controller)
			testCase.mockBehaviour(service, testCase.inputTransaction)

			handler := NewTransactionHandler(service, auth_service)
			router := mux.NewRouter()
			router.HandleFunc("/api/transactions/", handler.CreateTransaction)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/transactions/", bytes.NewBuffer(testCase.requestBody))

			router.ServeHTTP(w, r)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_CancelTransaction(t *testing.T) {
	// Arrange
	type mockBehaviour func(service *mock_services.MockTransactionService, transactionId int)
	// serializedTransaction, _ := json.Marshal(transaction)
	cancelSerializedTransaction, _ := json.Marshal(canceledTransaction)

	testTable := []struct {
		name                string
		transactionId       int
		requestBody         []byte
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Test cancel transaction (ok)",
			transactionId:       transaction.Id,
			expectedStatusCode:  http.StatusOK,
			requestBody:         emptyBody,
			expectedRequestBody: string(cancelSerializedTransaction) + "\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, transactionId int) {
				service.EXPECT().UpdateStatus(transactionId, services.TransactionCanceledStatus).Return(canceledTransaction, nil)
			},
		},
		{
			name:                "Test cancel transaction (not found)",
			transactionId:       transaction.Id,
			expectedStatusCode:  http.StatusNotFound,
			requestBody:         emptyBody,
			expectedRequestBody: "Not Found\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, transactionId int) {
				service.EXPECT().UpdateStatus(transactionId, services.TransactionCanceledStatus).Return(nil, errors.New(dbNotFoundErrorMsg))
			},
		},
		{
			name:                "Test cancel transaction (terminal status)",
			transactionId:       transaction.Id,
			expectedStatusCode:  http.StatusBadRequest,
			requestBody:         emptyBody,
			expectedRequestBody: "Can not proceed transaction with it's current status.\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, transactionId int) {
				service.EXPECT().UpdateStatus(transactionId, services.TransactionCanceledStatus).Return(nil, errors.New(services.TerminalStatusErrorMessage))
			},
		},
		{
			name:                "Test cancel transaction (service error)",
			transactionId:       transaction.Id,
			expectedStatusCode:  http.StatusInternalServerError,
			requestBody:         emptyBody,
			expectedRequestBody: "Internal Server Error\n",
			mockBehaviour: func(service *mock_services.MockTransactionService, transactionId int) {
				service.EXPECT().UpdateStatus(transactionId, services.TransactionCanceledStatus).Return(nil, errors.New("some error"))
			},
		},
	}

	// Act
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_services.NewMockTransactionService(controller)
			auth_service := mock_services.NewMockAuthMiddleware(controller)
			testCase.mockBehaviour(service, testCase.transactionId)

			handler := NewTransactionHandler(service, auth_service)
			router := mux.NewRouter()
			router.HandleFunc(fmt.Sprintf("/api/transactions/{pk:[0-9]}/cancel/"), handler.CancelTransaction)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", fmt.Sprintf("/api/transactions/%d/cancel/", testCase.transactionId), bytes.NewBuffer(testCase.requestBody))

			router.ServeHTTP(w, r)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
