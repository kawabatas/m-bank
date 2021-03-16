package model

import (
	"time"
)

type PaymentTransaction struct {
	UUID        string
	UserID      uint
	Amount      int // 負の数もとりうる
	TryTime     time.Time
	ConfirmTime time.Time
	CancelTime  time.Time
}

func NewPaymentTransaction(uuid string, userID uint, amount int) *PaymentTransaction {
	return &PaymentTransaction{
		UUID:    uuid,
		UserID:  userID,
		Amount:  amount,
		TryTime: time.Now(),
	}
}

func (pt *PaymentTransaction) IsTryStatus() bool {
	if pt.ConfirmTime.IsZero() && pt.CancelTime.IsZero() {
		return true
	}
	return false
}
