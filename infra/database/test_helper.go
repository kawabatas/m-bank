package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kawabatas/m-bank/domain/model"
)

func TestDBName() string {
	dbname := os.Getenv("DB_NAME")
	return fmt.Sprintf("%s_test", dbname)
}

func newTestConnection(t *testing.T) *sql.DB {
	db := newTestDBConnection(t)
	if err := truncateTables(db); err != nil {
		t.Fatal(err)
	}
	return db
}

var testDB *sql.DB
var initTestDB sync.Once

func newTestDBConnection(t *testing.T) *sql.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := TestDBName()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	// 初回のみDBの初期化を行う
	initTestDB.Do(func() {
		db, err := sql.Open("mysql", DSN(host, user, password, dbName))
		if err != nil {
			t.Fatal(err)
		}
		testDB = db
		testDB.SetMaxOpenConns(100)
		testDB.SetMaxIdleConns(100)
		testDB.SetConnMaxLifetime(10 * time.Second)
	})
	wg.Done()
	wg.Wait()

	return testDB
}

func truncateTables(db *sql.DB) error {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema IN (?) AND table_name != 'migrations'", TestDBName())
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		// FIXME: 以下の3つの文が同じコネクションで実行されるとは限らないが、大抵の場合はうまく動いている。問題が出てきたら修正する。
		// @see https://github.com/go-sql-driver/mysql/issues/373
		if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
			return err
		}
		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
		if _, err := db.Exec(truncateQuery); err != nil {
			return err
		}
		if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=1"); err != nil {
			return err
		}
	}
	return nil
}

const initBalanceAmount = 1000

func createSampleUsers(t *testing.T, db *sql.DB, count int) []*model.User {
	t.Helper()
	var users []*model.User
	var balances []*model.Balance
	for i := 1; i <= count; i++ {
		user := &model.User{
			ID:   uint(i),
			Name: fmt.Sprintf("sample%d", i),
		}
		balance := &model.Balance{
			UserID: uint(i),
			Amount: initBalanceAmount,
		}
		users = append(users, user)
		balances = append(balances, balance)
	}
	ctx := context.Background()
	var userStrings, balanceStrings []string
	var userArgs, balanceArgs []interface{}
	for _, u := range users {
		userStrings = append(userStrings, "(?, ?)")
		userArgs = append(userArgs, u.ID)
		userArgs = append(userArgs, u.Name)
	}
	for _, b := range balances {
		balanceStrings = append(balanceStrings, "(?, ?)")
		balanceArgs = append(balanceArgs, b.UserID)
		balanceArgs = append(balanceArgs, b.Amount)
	}
	userQuery := "INSERT INTO users (id, name) VALUES %s"
	userQuery = fmt.Sprintf(userQuery, strings.Join(userStrings, ","))
	if _, err := db.ExecContext(ctx, userQuery, userArgs...); err != nil {
		t.Fatalf("insert users error: %v", err)
	}
	balanceQuery := "INSERT INTO balances (user_id, amount) VALUES %s"
	balanceQuery = fmt.Sprintf(balanceQuery, strings.Join(balanceStrings, ","))
	if _, err := db.ExecContext(ctx, balanceQuery, balanceArgs...); err != nil {
		t.Fatalf("insert balances error: %v", err)
	}
	return users
}

func createSamplePaymentTransaction(t *testing.T, db *sql.DB, pt *model.PaymentTransaction) {
	t.Helper()
	ctx := context.Background()
	var confirmTime, cancelTime sql.NullTime
	if !pt.ConfirmTime.IsZero() {
		confirmTime.Valid = true
		confirmTime.Time = pt.ConfirmTime
	}
	if !pt.CancelTime.IsZero() {
		cancelTime.Valid = true
		cancelTime.Time = pt.CancelTime
	}
	if _, err := db.ExecContext(ctx,
		"INSERT INTO payment_transactions (uuid, user_id, amount, try_time, confirm_time, cancel_time) VALUES (?, ?, ?, ?, ?, ?)",
		pt.UUID, pt.UserID, pt.Amount, pt.TryTime, confirmTime, cancelTime,
	); err != nil {
		t.Fatalf("insert payment_transactions error: %v", err)
	}
}
