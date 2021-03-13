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

func NewPaymentTransaction(uuid string, userID uint, amount int) *PaymentTransaction {
	return &PaymentTransaction{
		UUID:       uuid,
		UserID:     userID,
		Amount:     amount,
		CreateTime: time.Now(),
	}
}
