package database

import (
	"context"
	"database/sql"

	"github.com/kawabatas/m-bank/domain/model"
)

type BalanceRepository struct {
	DB *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{DB: db}
}

func (r *BalanceRepository) Get(ctx context.Context, userID uint) (*model.Balance, error) {
	return &model.Balance{}, nil
}

func (r *BalanceRepository) AddAllUsers(ctx context.Context, amount, limit, offset int) error {
	return nil
}
