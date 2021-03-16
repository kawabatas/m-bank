package repository

import (
	"context"

	"github.com/kawabatas/m-bank/domain/model"
)

type PaymentTransactionRepository interface {
	Get(ctx context.Context, uuid string) (*model.PaymentTransaction, error)
	Try(ctx context.Context, uuid string, userID uint, amount int) (*model.PaymentTransaction, error)
	Confirm(ctx context.Context, uuid string) (*model.PaymentTransaction, error)
	Cancel(ctx context.Context, uuid string) (*model.PaymentTransaction, error)
}
