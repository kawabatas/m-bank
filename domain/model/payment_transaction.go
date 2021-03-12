package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentTransaction struct {
	UUID        string
	CreateTime  time.Time
	TryTime     *time.Time
	ConfirmTime *time.Time
	CancelTime  *time.Time
}

func NewPaymentTransaction() *PaymentTransaction {
	return &PaymentTransaction{
		UUID:       uuid.New().String(),
		CreateTime: time.Now(),
	}
}

func (t *PaymentTransaction) Try() {
	now := time.Now()
	t.TryTime = &now
}

func (t *PaymentTransaction) Confirm() {
	now := time.Now()
	t.ConfirmTime = &now
}

func (t *PaymentTransaction) Cancel() {
	now := time.Now()
	t.CancelTime = &now
}
