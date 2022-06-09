package repositories

import (
	"fmt"

	"github.com/Pythonyan3/payment-service/internal/database"
	"github.com/Pythonyan3/payment-service/internal/models"
)

var transactionTableName = "transaction"

type TransactionPostgresRepository struct {
	db *database.PostgresDB
}

func NewTransactionPostgresRepository(db *database.PostgresDB) *TransactionPostgresRepository {
	/*Transaction postgres repository constructor function.*/
	return &TransactionPostgresRepository{db: db}
}

func (repo *TransactionPostgresRepository) CreateTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	/*Insert new transaction data to DB and return transaction struct filled with new transaction data.*/

	// start new db transaction
	dbTransaction, err := repo.db.Beginx()

	if err != nil {
		return nil, err
	}

	// build query string
	query := fmt.Sprintf(
		"INSERT INTO %s (user_id, user_email, amount, currency, status) values ($1, $2, $3, $4, $5) RETURNING *",
		transactionTableName)

	// evalate insert query and parse new row data to transaction struct
	row := dbTransaction.QueryRowx(
		query, transaction.UserId, transaction.UserEmail, transaction.Amount, transaction.Currency, transaction.Status)
	err = row.StructScan(transaction)

	if err != nil {
		// roll back db transaction if parsing data to struct was failed
		dbTransaction.Rollback()
		return nil, err
	}

	// return transaction struct filled with data and commit db transaction
	return transaction, dbTransaction.Commit()
}

func (repo *TransactionPostgresRepository) UpdateTransactionStatus(transaction *models.Transaction, status string) (*models.Transaction, error) {
	/*Update transaction status return transaction struct filled with new transaction data.*/

	// start new db transaction
	dbTransaction, err := repo.db.Beginx()

	if err != nil {
		return nil, err
	}

	// build query string
	query := fmt.Sprintf(
		"UPDATE %s SET status = $1, updated_at = now()::timestamptz WHERE id = $2 RETURNING *;",
		transactionTableName)

	// evalate update query and parse new row data to transaction struct
	row := dbTransaction.QueryRowx(query, status, transaction.Id)
	err = row.StructScan(transaction)

	if err != nil {
		// roll back db transaction if parsing data to struct was failed
		dbTransaction.Rollback()
		return nil, err
	}

	// return transaction struct filled with data and commit db transaction
	return transaction, dbTransaction.Commit()
}

func (repo *TransactionPostgresRepository) GetTransactionById(transactionId int) (*models.Transaction, error) {
	/*Return transaction struct retrieved from db by PK.*/
	var transaction models.Transaction = models.Transaction{}

	// build query string
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", transactionTableName)

	// evaluate query and parse data to transaction struct
	if err := repo.db.Get(&transaction, query, transactionId); err != nil {
		return nil, err
	}

	return &transaction, nil
}
