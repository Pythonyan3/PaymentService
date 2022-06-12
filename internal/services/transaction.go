package services

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/Pythonyan3/payment-service/internal/models"
)

const (
	// Transaction Entity possible statuses
	TransactionNewStatus      string = "NEW"
	TransactionErrorStatus    string = "ERROR"
	TransactionFailedStatus   string = "FAILED"
	TransactionSuccessStatus  string = "SUCCESS"
	TransactionCanceledStatus string = "CANCELED"

	// Error message for updating transactions with terminal status
	TerminalStatusErrorMessage string = "TransactionService: cannot update transaction with it's current status."
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) (*models.Transaction, error)
	UpdateTransactionStatus(transaction *models.Transaction, status string) (*models.Transaction, error)
	GetTransactionById(transactionId int) (*models.Transaction, error)
}

type TransactionService struct {
	repo TransactionRepository
}

func NewTransactionService(repo TransactionRepository) *TransactionService {
	/*Transaction service constructor function.*/
	return &TransactionService{repo: repo}
}

func (service *TransactionService) Create(transactionInput *models.TransactionInput) (*models.Transaction, error) {
	/*Create new transaction (add new record to DB).*/
	var transaction models.Transaction = models.Transaction{
		UserId:    transactionInput.UserId,
		Amount:    transactionInput.Amount,
		Currency:  transactionInput.Currency,
		UserEmail: transactionInput.UserEmail,
	}

	// 1/5 of all trasnactions should be created with 'Error' status
	if rand.Intn(100) > 80 {
		transaction.Status = TransactionErrorStatus
	} else {
		transaction.Status = TransactionNewStatus
	}

	return service.repo.CreateTransaction(&transaction)
}

func (service *TransactionService) UpdateStatus(transactionId int, status string) (*models.Transaction, error) {
	/*Perform transaction status update (allowed only for transactions with status 'NEW').*/
	var transaction *models.Transaction
	var err error

	// retrieve transaction from db to update
	transaction, err = service.repo.GetTransactionById(transactionId)

	if err != nil {
		return nil, fmt.Errorf("service.repo.GetTransactionById failed: %w", err)
	}

	// check transaction current status and return error if it has not 'NEW' status
	if transaction.Status != TransactionNewStatus {
		return nil, errors.New(TerminalStatusErrorMessage)
	}

	// update transaction status
	transaction, err = service.repo.UpdateTransactionStatus(transaction, status)

	if err != nil {
		return nil, fmt.Errorf("service.repo.UpdateTransactionStatus failed: %w", err)
	}

	return transaction, nil
}

func (service *TransactionService) GetById(transactionId int) (*models.Transaction, error) {
	/*Retriving transaction data from DB by transaction PK.*/
	return service.repo.GetTransactionById(transactionId)
}
