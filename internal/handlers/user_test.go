package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pythonyan3/payment-service/internal/models"
	mock_services "github.com/Pythonyan3/payment-service/internal/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_TransactionsListByUserId(t *testing.T) {
	// Arrange
	type mockBehaviour func(service *mock_services.MockUserService, userId int)
	timeNow, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-06-12T18:09:14.796895+03:00")

	testTable := []struct {
		name                string
		userId              int
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Test list of transactions (ok)",
			userId:              1,
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "[{\"id\":1,\"user_id\":1,\"user_email\":\"email@mail.ru\",\"amount\":1500,\"currency\":\"EUR\",\"created_at\":\"2022-06-12T18:09:14.796895+03:00\",\"updated_at\":\"2022-06-12T18:09:14.796895+03:00\",\"status\":\"NEW\"}]\n",
			mockBehaviour: func(service *mock_services.MockUserService, userId int) {
				service.EXPECT().GetUserTransactionsById(userId).Return([]*models.Transaction{
					{
						Id:        1,
						UserId:    1,
						UserEmail: "email@mail.ru",
						Amount:    1500,
						Currency:  "EUR",
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
						Status:    "NEW",
					},
				}, nil)
			},
		},
		{
			name:                "Test list of transactions (empty)",
			userId:              1,
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "[]\n",
			mockBehaviour: func(service *mock_services.MockUserService, userId int) {
				service.EXPECT().GetUserTransactionsById(userId).Return([]*models.Transaction{}, nil)
			},
		},
		{
			name:                "Test list of transactions (service error)",
			userId:              1,
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
	timeNow, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-06-12T18:09:14.796895+03:00")

	testTable := []struct {
		name                string
		userEmail           string
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Test list of transactions (ok)",
			userEmail:           "email@mail.ru",
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "[{\"id\":1,\"user_id\":1,\"user_email\":\"email@mail.ru\",\"amount\":1500,\"currency\":\"EUR\",\"created_at\":\"2022-06-12T18:09:14.796895+03:00\",\"updated_at\":\"2022-06-12T18:09:14.796895+03:00\",\"status\":\"NEW\"}]\n",
			mockBehaviour: func(service *mock_services.MockUserService, userEmail string) {
				service.EXPECT().GetUserTransactionsByEmail(userEmail).Return([]*models.Transaction{
					{
						Id:        1,
						UserId:    1,
						UserEmail: "email@mail.ru",
						Amount:    1500,
						Currency:  "EUR",
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
						Status:    "NEW",
					},
				}, nil)
			},
		},
		{
			name:                "Test list of transactions (empty)",
			userEmail:           "email@mail.ru",
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "[]\n",
			mockBehaviour: func(service *mock_services.MockUserService, userEmail string) {
				service.EXPECT().GetUserTransactionsByEmail(userEmail).Return([]*models.Transaction{}, nil)
			},
		},
		{
			name:                "Test list of transactions (service error)",
			userEmail:           "email@mail.ru",
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
