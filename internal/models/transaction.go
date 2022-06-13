package models

import "time"

// base Transaction entity struct
type Transaction struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	UserEmail string    `json:"user_email" db:"user_email"`
	Amount    int64     `json:"amount" db:"amount"`
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Status    string    `json:"status" db:"status"`
}

// Transaction entity struct used for creating new transaction in API
// use validation tags for validation request data
type TransactionInput struct {
	UserId    int    `json:"user_id" validate:"required,gt=0"`
	UserEmail string `json:"user_email" validate:"required,email"`
	Amount    int64  `json:"amount" validate:"required"`
	Currency  string `json:"currency" validate:"required,len=3,uppercase"`
}

// Transaction status struct used for updating transaction status in API
// use validation tags for validation request data
type TransactionStatusInput struct {
	Status string `json:"status" validate:"required,uppercase,oneof=SUCCESS FAILED"`
}
