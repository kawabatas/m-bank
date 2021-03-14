package repository

import (
	"context"

	"github.com/kawabatas/m-bank/domain/model"
)

type BalanceRepository interface {
	Get(ctx context.Context, userID uint) (*model.Balance, error)
	AddToUsers(ctx context.Context, amount, limit, offset int) error
}
