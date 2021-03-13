package main

import (
	"database/sql"

	"github.com/kawabatas/m-bank/domain/repository"
	"github.com/kawabatas/m-bank/infra/database"
)

type application struct {
	BalanceService *balanceService
	PaymentService *paymentService
}

// newApp creates application services.
func newApp(db *sql.DB) *application {
	balanceRepository := database.NewBalanceRepository(db)
	paymentRepository := database.NewPaymentTransactionRepository(db)

	return &application{
		BalanceService: &balanceService{
			BalanceRepo: balanceRepository,
		},
		PaymentService: &paymentService{
			PaymentRepo: paymentRepository,
		},
	}
}

// balanceService is a service to handle balances.
type balanceService struct {
	BalanceRepo repository.BalanceRepository
}

// paymentService is a service to handle payments.
type paymentService struct {
	PaymentRepo repository.PaymentTransactionRepository
}
