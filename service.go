package main

import (
	"context"
	"database/sql"

	"github.com/kawabatas/m-bank/domain"
	"github.com/kawabatas/m-bank/domain/model"
	"github.com/kawabatas/m-bank/domain/repository"
	"github.com/kawabatas/m-bank/infra/database"
)

type application struct {
	BalanceService *balanceService
	PaymentService *paymentService
}

// balanceService is a service to handle balances.
type balanceService struct {
	BalanceRepo repository.BalanceRepository
}

// paymentService is a service to handle payments.
type paymentService struct {
	BalanceRepo repository.BalanceRepository
	PaymentRepo repository.PaymentTransactionRepository
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
			BalanceRepo: balanceRepository,
			PaymentRepo: paymentRepository,
		},
	}
}

func (s *balanceService) Get(ctx context.Context, userID uint) (*model.Balance, error) {
	return s.BalanceRepo.Get(ctx, userID)
}

func (s *paymentService) Try(ctx context.Context, uuid string, userID uint, amount int) (*model.PaymentTransaction, *model.Balance, error) {
	// TODO: 残高が足りるかチェック
	pt, err := s.PaymentRepo.Try(ctx, uuid, userID, amount)
	if err != nil {
		return nil, nil, err
	}

	balance, err := s.BalanceRepo.Get(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	return pt, balance, nil
}

func (s *paymentService) Confirm(ctx context.Context, uuid string, userID uint, amount int) (*model.PaymentTransaction, *model.Balance, error) {
	// TODO: 残高が足りるかチェック
	pt, err := s.PaymentRepo.Get(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}
	if !pt.IsTryStatus() {
		return nil, nil, domain.ErrInvalidUUID
	}
	pt, err = s.PaymentRepo.Confirm(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}

	balance, err := s.BalanceRepo.Get(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	return pt, balance, nil
}

func (s *paymentService) Cancel(ctx context.Context, uuid string, userID uint, amount int) (*model.PaymentTransaction, *model.Balance, error) {
	pt, err := s.PaymentRepo.Get(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}
	if !pt.IsTryStatus() {
		return nil, nil, domain.ErrInvalidUUID
	}
	pt, err = s.PaymentRepo.Cancel(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}

	balance, err := s.BalanceRepo.Get(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	return pt, balance, nil
}
