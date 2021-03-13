package repository

import (
	"context"

	"github.com/kawabatas/m-bank/domain/model"
)

type PaymentTransactionRepository interface {
	Find(ctx context.Context, UUID string) (*model.PaymentTransaction, error)
	Store(ctx context.Context, transaction *model.PaymentTransaction) (*model.PaymentTransaction, error)
}
