package database

import (
	"context"
	"database/sql"

	"github.com/kawabatas/m-bank/domain"
	"github.com/kawabatas/m-bank/domain/model"
)

type BalanceRepository struct {
	DB *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{DB: db}
}

func (r *BalanceRepository) Get(ctx context.Context, userID uint) (*model.Balance, error) {
	return findBalance(ctx, r.DB, userID)
}

func (r *BalanceRepository) AddAllUsers(ctx context.Context, amount, limit, offset int) error {
	return nil
}

func findBalance(ctx context.Context, db dbContext, userID uint) (*model.Balance, error) {
	query := `SELECT user_id, amount FROM balances WHERE user_id = ?`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, domain.ErrNoSuchEntity
	}
	return rowsToBalance(rows)
}

func rowsToBalance(rows *sql.Rows) (*model.Balance, error) {
	balance := &model.Balance{}
	if err := rows.Scan(&balance.UserID, &balance.Amount); err != nil {
		return nil, err
	}
	return balance, nil
}
