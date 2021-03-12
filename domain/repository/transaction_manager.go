package repository

import (
	"context"
)

type TransactionFunc func(ctx context.Context) error

type TransactionManager interface {
	WithTransaction(ctx context.Context, innerFunc TransactionFunc) (err error)
}

var manager TransactionManager

func SetTransactionManager(m TransactionManager) {
	manager = m
}

func WithTransaction(ctx context.Context, innerFunc TransactionFunc) (err error) {
	return manager.WithTransaction(ctx, innerFunc)
}
