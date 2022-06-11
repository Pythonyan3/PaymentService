package repositories

import (
	"fmt"

	"github.com/Pythonyan3/payment-service/internal/database"
	"github.com/Pythonyan3/payment-service/internal/models"
)

type UserPostgresRepository struct {
	db *database.PostgresDB
}

func NewUserPostgresRepository(db *database.PostgresDB) *UserPostgresRepository {
	/*User postgres repository constructor function.*/
	return &UserPostgresRepository{db: db}
}

func (repo *UserPostgresRepository) GetUserTransactionsById(userId int) ([]*models.Transaction, error) {
	/*Return slice of transaction structs retrieved from db filtered by user id.*/
	var transactions []*models.Transaction = make([]*models.Transaction, 0)
	var query string

	// build query string
	query = fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY created_at DESC;", transactionTableName)

	// evaluate query and parse data to slice of transaction structs
	if err := repo.db.Select(&transactions, query, userId); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *UserPostgresRepository) GetUserTransactionsByEmail(userEmail string) ([]*models.Transaction, error) {
	/*Return slice of transaction structs retrieved from db filtered by user email.*/
	var transactions []*models.Transaction = make([]*models.Transaction, 0)
	var query string

	// build query string
	query = fmt.Sprintf("SELECT * FROM %s WHERE user_email = $1 ORDER BY created_at DESC;", transactionTableName)

	// evaluate query and parse data to slice of transaction structs
	if err := repo.db.Select(&transactions, query, userEmail); err != nil {
		return nil, err
	}

	return transactions, nil
}
