package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

func (r *BalanceRepository) AddToUsers(ctx context.Context, amount, limit, offset int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 対象のユーザの残高取得
	var balances []model.Balance
	fetchQuery := `SELECT user_id, amount FROM balances ORDER BY user_id ASC LIMIT ? OFFSET ? FOR UPDATE`
	rows, err := tx.QueryContext(ctx, fetchQuery, limit, offset)
	if err != nil {
		return err
	}
	for rows.Next() {
		var balance model.Balance
		if err := rows.Scan(&balance.UserID, &balance.Amount); err != nil {
			return err
		}
		balances = append(balances, balance)
	}

	// 残高加算
	ids := make([]interface{}, len(balances))
	for i, b := range balances {
		ids[i] = b.UserID
	}
	updateQuery := "UPDATE balances SET amount = amount + " + strconv.Itoa(amount) + " WHERE user_id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"
	if _, err := tx.ExecContext(ctx, updateQuery, ids...); err != nil {
		return err
	}

	// ログ挿入(bulk insert)
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, b := range balances {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, b.UserID)
		valueArgs = append(valueArgs, b.Amount)
		valueArgs = append(valueArgs, b.Amount+uint(amount))
	}
	insertQuery := "INSERT INTO balance_logs (user_id, before_amount, after_amount) VALUES %s"
	insertQuery = fmt.Sprintf(insertQuery, strings.Join(valueStrings, ","))
	if _, err := tx.ExecContext(ctx, insertQuery, valueArgs...); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
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

func countBalanceLog(ctx context.Context, db dbContext, userID uint) (int, error) {
	query := `SELECT COUNT(id) FROM balance_logs WHERE user_id = ?`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, nil
	}
	var count int
	if err := rows.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func rowsToBalance(rows *sql.Rows) (*model.Balance, error) {
	balance := &model.Balance{}
	if err := rows.Scan(&balance.UserID, &balance.Amount); err != nil {
		return nil, err
	}
	return balance, nil
}
