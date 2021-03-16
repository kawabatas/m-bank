package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kawabatas/m-bank/domain/model"
)

func newPaymentTransactionRepo(t *testing.T) *PaymentTransactionRepository {
	t.Helper()
	db := newTestConnection(t)
	return NewPaymentTransactionRepository(db)
}

func TestPaymentTransactionRepository_Get(t *testing.T) {
	repo := newPaymentTransactionRepo(t)
	users := createSampleUsers(t, repo.DB, 1)
	sampleUuid := "foo"
	sampleAmount := 100
	createSamplePaymentTransaction(t, repo.DB,
		&model.PaymentTransaction{
			UUID:    sampleUuid,
			UserID:  users[0].ID,
			Amount:  sampleAmount,
			TryTime: time.Now(),
		},
	)
	ctx := context.Background()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.PaymentTransaction
		wantErr bool
	}{
		{
			"取得できる",
			fields{repo.DB},
			args{ctx, sampleUuid},
			&model.PaymentTransaction{
				UUID:   sampleUuid,
				UserID: users[0].ID,
				Amount: sampleAmount,
			},
			false,
		},
		{
			"取得できない",
			fields{repo.DB},
			args{ctx, "wrong uuid"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PaymentTransactionRepository{
				DB: tt.fields.DB,
			}
			got, err := r.Get(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentTransactionRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(model.PaymentTransaction{}, "TryTime", "ConfirmTime", "CancelTime")
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("PaymentTransactionRepository.Get() mismatch (-want +got): \n %s", diff)
			}
		})
	}
}

func TestPaymentTransactionRepository_Try(t *testing.T) {
	repo := newPaymentTransactionRepo(t)
	users := createSampleUsers(t, repo.DB, 2)
	sampleUuid := "foo"
	sampleAmount := 100
	ctx := context.Background()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx    context.Context
		uuid   string
		userID uint
		amount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.PaymentTransaction
		wantErr bool
	}{
		{
			"作成できる",
			fields{repo.DB},
			args{ctx, sampleUuid, users[0].ID, sampleAmount},
			&model.PaymentTransaction{
				UUID:   sampleUuid,
				UserID: users[0].ID,
				Amount: sampleAmount,
			},
			false,
		},
		{
			"同じUUIDでは作成できない",
			fields{repo.DB},
			args{ctx, sampleUuid, users[1].ID, sampleAmount},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PaymentTransactionRepository{
				DB: tt.fields.DB,
			}
			got, err := r.Try(tt.args.ctx, tt.args.uuid, tt.args.userID, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentTransactionRepository.Try() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(model.PaymentTransaction{}, "TryTime", "ConfirmTime", "CancelTime")
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("PaymentTransactionRepository.Try() mismatch (-want +got): \n %s", diff)
			}
			if got != nil {
				if got.TryTime.IsZero() {
					t.Error("PaymentTransactionRepository.Try() got.TryTime MUST NOT IsZero")
				}
				if !got.ConfirmTime.IsZero() {
					t.Error("PaymentTransactionRepository.Try() got.ConfirmTime MUST IsZero")
				}
				if !got.CancelTime.IsZero() {
					t.Error("PaymentTransactionRepository.Try() got.CancelTime MUST IsZero")
				}
			}
		})
	}
}

func TestPaymentTransactionRepository_Confirm(t *testing.T) {
	repo := newPaymentTransactionRepo(t)
	users := createSampleUsers(t, repo.DB, 1)
	addUuid := "add"
	notTryUuid := "not try"
	subUuid := "sub"
	subUuid2 := "sub minus"
	createSamplePaymentTransaction(t, repo.DB,
		&model.PaymentTransaction{
			UUID:    addUuid,
			UserID:  users[0].ID,
			Amount:  1,
			TryTime: time.Now(),
		},
	)
	createSamplePaymentTransaction(t, repo.DB,
		&model.PaymentTransaction{
			UUID:        notTryUuid,
			UserID:      users[0].ID,
			Amount:      1,
			TryTime:     time.Now(),
			ConfirmTime: time.Now(),
		},
	)
	createSamplePaymentTransaction(t, repo.DB,
		&model.PaymentTransaction{
			UUID:    subUuid,
			UserID:  users[0].ID,
			Amount:  -(initBalanceAmount + 1),
			TryTime: time.Now(),
		},
	)
	createSamplePaymentTransaction(t, repo.DB,
		&model.PaymentTransaction{
			UUID:    subUuid2,
			UserID:  users[0].ID,
			Amount:  -1,
			TryTime: time.Now(),
		},
	)
	ctx := context.Background()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		uuid string
	}
	type want2 struct {
		Balance  uint
		LogCount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.PaymentTransaction
		want2   want2
		wantErr bool
	}{
		{
			"加算できる",
			fields{repo.DB},
			args{ctx, addUuid},
			&model.PaymentTransaction{
				UUID:   addUuid,
				UserID: users[0].ID,
				Amount: 1,
			},
			want2{initBalanceAmount + 1, 1},
			false,
		},
		{
			"ステータスがtryでないUUID",
			fields{repo.DB},
			args{ctx, notTryUuid},
			nil,
			want2{initBalanceAmount + 1, 1},
			true,
		},
		{
			"減算できる",
			fields{repo.DB},
			args{ctx, subUuid},
			&model.PaymentTransaction{
				UUID:   subUuid,
				UserID: users[0].ID,
				Amount: -(initBalanceAmount + 1),
			},
			want2{0, 2},
			false,
		},
		{
			"残高がマイナスになる減算はできない",
			fields{repo.DB},
			args{ctx, subUuid2},
			nil,
			want2{0, 2},
			true,
		},
		{
			"存在しないUUID",
			fields{repo.DB},
			args{ctx, "wrong"},
			nil,
			want2{0, 2},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PaymentTransactionRepository{
				DB: tt.fields.DB,
			}
			got, err := r.Confirm(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentTransactionRepository.Confirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(model.PaymentTransaction{}, "TryTime", "ConfirmTime", "CancelTime")
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("PaymentTransactionRepository.Confirm() mismatch (-want +got): \n %s", diff)
			}
			if got != nil {
				if got.ConfirmTime.IsZero() {
					t.Error("PaymentTransactionRepository.Confirm() got.ConfirmTime MUST NOT IsZero")
				}
				b, err := findBalance(context.Background(), r.DB, got.UserID)
				if err != nil {
					t.Errorf("PaymentTransactionRepository.Confirm() findBalance error = %v", err)
				}
				c, err := countBalanceLog(context.Background(), r.DB, got.UserID)
				if err != nil {
					t.Errorf("PaymentTransactionRepository.Confirm() countBalanceLog error = %v", err)
				}
				got2 := want2{Balance: b.Amount, LogCount: c}
				if diff2 := cmp.Diff(tt.want2, got2, nil); diff2 != "" {
					t.Errorf("PaymentTransactionRepository.Confirm() mismatch (-want2 +got2): \n %s", diff2)
				}
			}
		})
	}
}

func TestPaymentTransactionRepository_Cancel(t *testing.T) {
	repo := newPaymentTransactionRepo(t)
	users := createSampleUsers(t, repo.DB, 1)
	tryUuid := "try"
	notTryUuid := "not try"
	sampleAmount := 100
	createSamplePaymentTransaction(t, repo.DB,
		&model.PaymentTransaction{
			UUID:    tryUuid,
			UserID:  users[0].ID,
			Amount:  sampleAmount,
			TryTime: time.Now(),
		},
	)
	createSamplePaymentTransaction(t, repo.DB,
		&model.PaymentTransaction{
			UUID:       notTryUuid,
			UserID:     users[0].ID,
			Amount:     sampleAmount,
			TryTime:    time.Now(),
			CancelTime: time.Now(),
		},
	)
	ctx := context.Background()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.PaymentTransaction
		wantErr bool
	}{
		{
			"キャンセルできる",
			fields{repo.DB},
			args{ctx, tryUuid},
			&model.PaymentTransaction{
				UUID:   tryUuid,
				UserID: users[0].ID,
				Amount: sampleAmount,
			},
			false,
		},
		{
			"ステータスがtryでないUUID",
			fields{repo.DB},
			args{ctx, notTryUuid},
			nil,
			true,
		},
		{
			"存在しないUUID",
			fields{repo.DB},
			args{ctx, "wrong"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PaymentTransactionRepository{
				DB: tt.fields.DB,
			}
			got, err := r.Cancel(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentTransactionRepository.Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(model.PaymentTransaction{}, "TryTime", "ConfirmTime", "CancelTime")
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("PaymentTransactionRepository.Cancel() mismatch (-want +got): \n %s", diff)
			}
			if got != nil {
				if got.CancelTime.IsZero() {
					t.Error("PaymentTransactionRepository.Cancel() got.CancelTime MUST NOT IsZero")
				}
			}
		})
	}
}
