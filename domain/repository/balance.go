package repository

import (
	"context"

	"github.com/kawabatas/m-bank/domain/model"
)

type BalanceRepository interface {
	Get(ctx context.Context, userID uint) (*model.Balance, error)
	Add(ctx context.Context, userID uint, amount int) (*model.Balance, error)
	AddAllUsers(ctx context.Context, amount int) error
}
