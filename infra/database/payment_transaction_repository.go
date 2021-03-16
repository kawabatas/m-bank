package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kawabatas/m-bank/domain"
	"github.com/kawabatas/m-bank/domain/model"
)

type PaymentTransactionRepository struct {
	DB *sql.DB
}

func NewPaymentTransactionRepository(db *sql.DB) *PaymentTransactionRepository {
	return &PaymentTransactionRepository{DB: db}
}

func (r *PaymentTransactionRepository) Get(ctx context.Context, uuid string) (*model.PaymentTransaction, error) {
	return findPaymentTransaction(ctx, r.DB, uuid, false)
}

func (r *PaymentTransactionRepository) Try(ctx context.Context, uuid string, userID uint, amount int) (*model.PaymentTransaction, error) {
	pt := model.NewPaymentTransaction(uuid, userID, amount)
	if _, err := r.DB.ExecContext(ctx,
		"INSERT INTO payment_transactions (uuid, user_id, amount, try_time) VALUES (?, ?, ?, ?)",
		pt.UUID, pt.UserID, pt.Amount, pt.TryTime,
	); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, domain.ErrDuplicateEntity
		}
		return nil, err
	}
	return findPaymentTransaction(ctx, r.DB, uuid, false)
}

func (r *PaymentTransactionRepository) Confirm(ctx context.Context, uuid string) (*model.PaymentTransaction, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	pt, err := findPaymentTransaction(ctx, tx, uuid, true)
	if err != nil {
		return nil, err
	}
	if !pt.IsTryStatus() {
		return nil, domain.ErrInvalidUUID
	}
	pt.ConfirmTime = time.Now()
	if _, err := tx.ExecContext(ctx,
		`UPDATE payment_transactions SET confirm_time = ? WHERE uuid = ?`,
		pt.ConfirmTime, pt.UUID,
	); err != nil {
		return nil, err
	}

	// トランザクション内で、残高不足のチェックおよび加減算を行う
	beforeBalance, err := findBalance(ctx, tx, pt.UserID)
	if err != nil {
		return nil, err
	}
	// 残高が負の値でないことはamountの型で担保
	if _, err := tx.ExecContext(ctx,
		`UPDATE balances SET amount = amount + ? WHERE user_id = ?`,
		pt.Amount, pt.UserID,
	); err != nil {
		return nil, err
	}
	if _, err := r.DB.ExecContext(ctx,
		"INSERT INTO balance_logs (user_id, before_amount, after_amount) VALUES (?, ?, ?)",
		pt.UserID, beforeBalance.Amount, beforeBalance.Amount+uint(pt.Amount),
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 再取得
	return findPaymentTransaction(ctx, r.DB, uuid, false)
}

func (r *PaymentTransactionRepository) Cancel(ctx context.Context, uuid string) (*model.PaymentTransaction, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	pt, err := findPaymentTransaction(ctx, tx, uuid, true)
	if err != nil {
		return nil, err
	}
	if !pt.IsTryStatus() {
		return nil, domain.ErrInvalidUUID
	}
	pt.CancelTime = time.Now()
	if _, err := tx.ExecContext(ctx,
		`UPDATE payment_transactions SET cancel_time = ? WHERE uuid = ?`,
		pt.CancelTime, pt.UUID,
	); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 再取得
	return findPaymentTransaction(ctx, r.DB, uuid, false)
}

func findPaymentTransaction(ctx context.Context, db dbContext, uuid string, withLock bool) (*model.PaymentTransaction, error) {
	query := `
	SELECT
		uuid, user_id, amount,
		try_time, confirm_time, cancel_time
	FROM payment_transactions WHERE uuid = ?`
	if withLock {
		query = query + ` FOR UPDATE`
	}
	rows, err := db.QueryContext(ctx, query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, domain.ErrNoSuchEntity
	}
	return rowsToPaymentTransaction(rows)
}

func rowsToPaymentTransaction(rows *sql.Rows) (*model.PaymentTransaction, error) {
	pt := &model.PaymentTransaction{}
	var confirmTime, cancelTime sql.NullTime
	if err := rows.Scan(&pt.UUID, &pt.UserID, &pt.Amount, &pt.TryTime, &confirmTime, &cancelTime); err != nil {
		return nil, err
	}
	if confirmTime.Valid {
		pt.ConfirmTime = confirmTime.Time
	}
	if cancelTime.Valid {
		pt.CancelTime = cancelTime.Time
	}
	return pt, nil
}
