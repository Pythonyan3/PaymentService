package services

import "github.com/Pythonyan3/payment-service/internal/models"

type UserRepository interface {
	GetUserTransactionsById(userId int) ([]*models.Transaction, error)
	GetUserTransactionsByEmail(userEmail string) ([]*models.Transaction, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	/*User service constructor function.*/
	return &UserService{repo: repo}
}

func (service *UserService) GetUserTransactionsById(userId int) ([]*models.Transaction, error) {
	/*Retrieve list of transactions filtered by user id.*/
	return service.repo.GetUserTransactionsById(userId)
}

func (service *UserService) GetUserTransactionsByEmail(userEmail string) ([]*models.Transaction, error) {
	/*Retrieve list of transactions filtered by user email.*/
	return service.repo.GetUserTransactionsByEmail(userEmail)
}
