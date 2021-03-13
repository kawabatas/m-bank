package database

import (
	"context"
	"database/sql"

	"github.com/kawabatas/m-bank/domain/model"
)

type PaymentTransactionRepository struct {
	DB *sql.DB
}

func NewPaymentTransactionRepository(db *sql.DB) *PaymentTransactionRepository {
	return &PaymentTransactionRepository{DB: db}
}

func (r *PaymentTransactionRepository) Find(ctx context.Context, UUID string) (*model.PaymentTransaction, error) {
	return nil, nil
}

func (r *PaymentTransactionRepository) Store(ctx context.Context, transaction *model.PaymentTransaction) (*model.PaymentTransaction, error) {
	return nil, nil
}
