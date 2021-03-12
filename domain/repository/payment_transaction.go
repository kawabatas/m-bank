package repository

import (
	"context"

	"github.com/kawabatas/m-bank/domain/model"
)

type CoinTransactionRepository interface {
	Find(ctx context.Context, UUID string) (*model.PaymentTransaction, error)
	CreateOrUpdate(ctx context.Context, transaction *model.PaymentTransaction) (*model.PaymentTransaction, error)
}
