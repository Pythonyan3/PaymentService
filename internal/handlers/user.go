package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Pythonyan3/payment-service/internal/models"

	"github.com/gorilla/mux"
)

type UserService interface {
	GetUserTransactionsById(userId int) ([]*models.Transaction, error)
	GetUserTransactionsByEmail(userEmail string) ([]*models.Transaction, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	/*User routes handler constructor function.*/
	return &UserHandler{service: service}
}

func (handler *UserHandler) InitRoutes(router *mux.Router) {
	/*Perform initialization of all required routes for user entity.*/
	var subRouter *mux.Router = router.PathPrefix("/users").Subrouter()
	subRouter.HandleFunc("/{userId:[0-9]+}/transactions/", handler.TransactionsListByUserId).Methods("GET")
	subRouter.HandleFunc("/{userEmail}/transactions/", handler.TransactionsListByUserEmail).Methods("GET")
}

func (handler *UserHandler) TransactionsListByUserId(w http.ResponseWriter, r *http.Request) {
	/*
		Handle request to retrieve list of user's transactions.

		Accept user PK in URL params to perform filtering.
	*/
	// declare variables
	var userId int
	var err error
	var transaction []*models.Transaction
	var params map[string]string = mux.Vars(r)

	// retrieve user PK from url variables
	userId, err = strconv.Atoi(params["userId"])
	if err != nil {
		log.Printf("strconv.Atoi failed: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// use service to retrieve user's transactions
	transaction, err = handler.service.GetUserTransactionsById(userId)

	if err != nil {
		log.Printf("handler.service.GetUserTransactionsById failed: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (handler *UserHandler) TransactionsListByUserEmail(w http.ResponseWriter, r *http.Request) {
	/*
		Handle request to retrieve list of user's transactions.

		Accept user email address in URL params to perform filtering.
	*/
	// declare variables
	var err error
	var transaction []*models.Transaction
	var params map[string]string = mux.Vars(r)

	// use service to retrieve user's transactions
	transaction, err = handler.service.GetUserTransactionsByEmail(params["userEmail"])

	if err != nil {
		log.Printf("handler.service.GetUserTransactionsByEmail failed: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}
