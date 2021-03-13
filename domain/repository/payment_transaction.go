package repository

import (
	"context"

	"github.com/kawabatas/m-bank/domain/model"
)

type PaymentTransactionRepository interface {
	Find(ctx context.Context, UUID string) (*model.PaymentTransaction, error)
	Try(ctx context.Context, UUID string) (*model.PaymentTransaction, error)
	Confirm(ctx context.Context, UUID string) (*model.PaymentTransaction, error)
	Cancel(ctx context.Context, UUID string) (*model.PaymentTransaction, error)
}
