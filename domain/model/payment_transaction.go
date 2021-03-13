package model

import (
	"time"
)

type PaymentTransaction struct {
	UUID        string
	UserID      uint
	Amount      int // 負の数もとりうる
	CreateTime  time.Time
	TryTime     *time.Time
	ConfirmTime *time.Time
	CancelTime  *time.Time
}

func NewPaymentTransaction(uuid string) *PaymentTransaction {
	return &PaymentTransaction{
		UUID:       uuid,
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
