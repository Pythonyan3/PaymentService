package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pythonyan3/payment-service/internal/models"
	"github.com/Pythonyan3/payment-service/internal/services"
	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
)

const dbNotFoundErrorMsg = "sql: no rows in result set"

type TransactionService interface {
	GetById(transactionId int) (*models.Transaction, error)
	Create(transactionInput *models.TransactionInput) (*models.Transaction, error)
	UpdateStatus(transactionId int, status string) (*models.Transaction, error)
}

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	/*Transaction routes handler constructor function.*/
	return &TransactionHandler{service: service}
}

func (handler *TransactionHandler) InitRoutes(router *mux.Router) {
	/*Perform initialization of all required routes for transaction entity.*/
	var subRouter *mux.Router = router.PathPrefix("/transactions").Subrouter()

	subRouter.HandleFunc("/", handler.CreateTransaction).Methods("POST")
	subRouter.HandleFunc("/{pk:[0-9]+}/", handler.RetrieveTransaction).Methods("GET")
	subRouter.HandleFunc("/{pk:[0-9]+}/cancel/", handler.CancelTransaction).Methods("PUT", "PATCH")
	subRouter.HandleFunc("/{pk:[0-9]+}/proceed/", handler.ProceedTransaction).Methods("PUT", "PATCH")
}

func (handler *TransactionHandler) RetrieveTransaction(w http.ResponseWriter, r *http.Request) {
	/*
		Handle request to retrieve transaction info.

		Accept transaction PK in URL params to perform retriving.
	*/
	var err error
	var transactionId int
	var transaction *models.Transaction
	var params map[string]string = mux.Vars(r)

	// retrieve transaction PK from url variables
	transactionId, err = strconv.Atoi(params["pk"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// retriving transaction info with service
	transaction, err = handler.service.GetById(transactionId)

	if err != nil {
		if strings.Contains(err.Error(), dbNotFoundErrorMsg) {
			// return HTTP 404 status code if transaction was not found
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			// otherwise probably something went wrong...
			log.Printf("handler.service.GetById failed: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (handler *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	/*Handle request to create new transaction.*/
	var err error
	var transactionInput models.TransactionInput
	var transaction *models.Transaction
	var validator *validator.Validate = validator.New()

	// parsing reqeust body data to transaction struct
	if err := json.NewDecoder(r.Body).Decode(&transactionInput); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest)
		return
	}

	// validate parsed data
	if err := validator.Struct(transactionInput); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest)
		return
	}

	// create new transaction with a service
	transaction, err = handler.service.Create(&transactionInput)
	if err != nil {
		log.Printf("handler.service.Create failed: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func (handler *TransactionHandler) ProceedTransaction(w http.ResponseWriter, r *http.Request) {
	/*Handle request to update transaction status by payment service.*/
	var err error
	var transactionId int
	var transaction *models.Transaction
	var transactionStatusInput models.TransactionStatusInput
	var params map[string]string = mux.Vars(r)
	var validator *validator.Validate = validator.New()

	// retrieve transaction PK from URL variables
	transactionId, err = strconv.Atoi(params["pk"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// parse request body data to transaction status struct
	if err := json.NewDecoder(r.Body).Decode(&transactionStatusInput); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest)
		return
	}

	// validate parsed data
	if err := validator.Struct(transactionStatusInput); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %s", err), http.StatusBadRequest)
		return
	}

	// update transaction status with a service struct
	transaction, err = handler.service.UpdateStatus(transactionId, transactionStatusInput.Status)

	if err != nil {
		if strings.Contains(err.Error(), services.TerminalStatusErrorMessage) {
			// return HTTP 400 status code if transaction has terminal status
			http.Error(w, "Can not proceed transaction with it's current status.", http.StatusBadRequest)
		} else if strings.Contains(err.Error(), dbNotFoundErrorMsg) {
			// return HTTP 404 status code if transaction was not found
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			// otherwise probably something went wrong...
			log.Printf("handler.service.UpdateStatus failed: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (handler *TransactionHandler) CancelTransaction(w http.ResponseWriter, r *http.Request) {
	/*Handle request to cancel transaction.*/
	var err error
	var transactionId int
	var transaction *models.Transaction
	var params map[string]string = mux.Vars(r)

	// retrieve transaction PK from URL variables
	transactionId, err = strconv.Atoi(params["pk"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// update transaction status with a service
	transaction, err = handler.service.UpdateStatus(transactionId, services.TransactionCanceledStatus)

	if err != nil {
		if strings.Contains(err.Error(), services.TerminalStatusErrorMessage) {
			// return HTTP 400 status code if transaction has terminal status
			http.Error(w, "Can not proceed transaction with it's current status.", http.StatusBadRequest)
		} else if strings.Contains(err.Error(), dbNotFoundErrorMsg) {
			// return HTTP 404 status code if transaction was not found
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			// otherwise probably something went wrong...
			log.Printf("handler.service.UpdateStatus failed: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(transaction)
}
