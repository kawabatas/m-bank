package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kawabatas/m-bank/domain"
	"github.com/kawabatas/m-bank/domain/mock"
	"github.com/kawabatas/m-bank/domain/model"
	"github.com/kawabatas/m-bank/domain/repository"
)

func Test_paymentService_Try(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sampleBalance := &model.Balance{
		UserID: 1,
		Amount: 100,
	}
	samplePayment := &model.PaymentTransaction{
		UUID:   "foo",
		UserID: 1,
		Amount: 100,
	}
	balanceRepo := mock.NewMockBalanceRepository(ctrl)
	balanceRepo.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(sampleBalance, nil).
		AnyTimes()
	paymentRepo := mock.NewMockPaymentTransactionRepository(ctrl)
	paymentRepo.
		EXPECT().
		Try(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(samplePayment, nil).
		Times(2)

	ctx := context.Background()

	type fields struct {
		BalanceRepo repository.BalanceRepository
		PaymentRepo repository.PaymentTransactionRepository
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
		want1   *model.Balance
		wantErr bool
	}{
		{
			"加算できる",
			fields{balanceRepo, paymentRepo},
			args{ctx, samplePayment.UUID, samplePayment.UserID, 1},
			samplePayment,
			sampleBalance,
			false,
		},
		{
			"減算できる",
			fields{balanceRepo, paymentRepo},
			args{ctx, samplePayment.UUID, samplePayment.UserID, -samplePayment.Amount},
			samplePayment,
			sampleBalance,
			false,
		},
		{
			"減算で残高が足りない",
			fields{balanceRepo, paymentRepo},
			args{ctx, samplePayment.UUID, samplePayment.UserID, -int(sampleBalance.Amount) - 1},
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				BalanceRepo: tt.fields.BalanceRepo,
				PaymentRepo: tt.fields.PaymentRepo,
			}
			got, got1, err := s.Try(tt.args.ctx, tt.args.uuid, tt.args.userID, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("paymentService.Try() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentService.Try() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("paymentService.Try() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_paymentService_Confirm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sampleBalance := &model.Balance{
		UserID: 1,
		Amount: 100,
	}
	invalidUuid := "invalid"
	notTryUuid := "not try"
	samplePayment := &model.PaymentTransaction{
		UUID:   "foo",
		UserID: 1,
		Amount: 100,
	}
	balanceRepo := mock.NewMockBalanceRepository(ctrl)
	balanceRepo.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(sampleBalance, nil).
		AnyTimes()
	paymentRepo := mock.NewMockPaymentTransactionRepository(ctrl)
	paymentRepo.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, uuid string) (*model.PaymentTransaction, error) {
			if uuid == invalidUuid || uuid == notTryUuid {
				return nil, domain.ErrInvalidUUID
			}
			return samplePayment, nil
		}).
		AnyTimes()
	paymentRepo.
		EXPECT().
		Confirm(gomock.Any(), gomock.Any()).
		Return(samplePayment, nil).
		Times(2)

	ctx := context.Background()

	type fields struct {
		BalanceRepo repository.BalanceRepository
		PaymentRepo repository.PaymentTransactionRepository
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
		want1   *model.Balance
		wantErr bool
	}{
		{
			"加算できる",
			fields{balanceRepo, paymentRepo},
			args{ctx, samplePayment.UUID, samplePayment.UserID, 1},
			samplePayment,
			sampleBalance,
			false,
		},
		{
			"減算できる",
			fields{balanceRepo, paymentRepo},
			args{ctx, samplePayment.UUID, samplePayment.UserID, -samplePayment.Amount},
			samplePayment,
			sampleBalance,
			false,
		},
		{
			"減算で残高が足りない",
			fields{balanceRepo, paymentRepo},
			args{ctx, samplePayment.UUID, samplePayment.UserID, -int(sampleBalance.Amount) - 1},
			nil,
			nil,
			true,
		},
		{
			"ステータスがtryでないUUID",
			fields{balanceRepo, paymentRepo},
			args{ctx, notTryUuid, samplePayment.UserID, 1},
			nil,
			nil,
			true,
		},
		{
			"存在しないUUID",
			fields{balanceRepo, paymentRepo},
			args{ctx, invalidUuid, samplePayment.UserID, 1},
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				BalanceRepo: tt.fields.BalanceRepo,
				PaymentRepo: tt.fields.PaymentRepo,
			}
			got, got1, err := s.Confirm(tt.args.ctx, tt.args.uuid, tt.args.userID, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("paymentService.Confirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentService.Confirm() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("paymentService.Confirm() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_paymentService_Cancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sampleBalance := &model.Balance{
		UserID: 1,
		Amount: 100,
	}
	invalidUuid := "invalid"
	notTryUuid := "not try"
	samplePayment := &model.PaymentTransaction{
		UUID:   "foo",
		UserID: 1,
		Amount: 100,
	}
	balanceRepo := mock.NewMockBalanceRepository(ctrl)
	balanceRepo.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(sampleBalance, nil).
		AnyTimes()
	paymentRepo := mock.NewMockPaymentTransactionRepository(ctrl)
	paymentRepo.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, uuid string) (*model.PaymentTransaction, error) {
			if uuid == invalidUuid || uuid == notTryUuid {
				return nil, domain.ErrInvalidUUID
			}
			return samplePayment, nil
		}).
		AnyTimes()
	paymentRepo.
		EXPECT().
		Cancel(gomock.Any(), gomock.Any()).
		Return(samplePayment, nil).
		Times(1)

	ctx := context.Background()

	type fields struct {
		BalanceRepo repository.BalanceRepository
		PaymentRepo repository.PaymentTransactionRepository
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
		want1   *model.Balance
		wantErr bool
	}{
		{
			"キャンセルできる",
			fields{balanceRepo, paymentRepo},
			args{ctx, samplePayment.UUID, samplePayment.UserID, 1},
			samplePayment,
			sampleBalance,
			false,
		},
		{
			"ステータスがtryでないUUID",
			fields{balanceRepo, paymentRepo},
			args{ctx, notTryUuid, samplePayment.UserID, 1},
			nil,
			nil,
			true,
		},
		{
			"存在しないUUID",
			fields{balanceRepo, paymentRepo},
			args{ctx, invalidUuid, samplePayment.UserID, 1},
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				BalanceRepo: tt.fields.BalanceRepo,
				PaymentRepo: tt.fields.PaymentRepo,
			}
			got, got1, err := s.Cancel(tt.args.ctx, tt.args.uuid, tt.args.userID, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("paymentService.Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentService.Cancel() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("paymentService.Cancel() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_paymentService_AddToUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balanceRepo := mock.NewMockBalanceRepository(ctrl)
	balanceRepo.
		EXPECT().
		AddToUsers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	paymentRepo := mock.NewMockPaymentTransactionRepository(ctrl)

	ctx := context.Background()

	type fields struct {
		BalanceRepo repository.BalanceRepository
		PaymentRepo repository.PaymentTransactionRepository
	}
	type args struct {
		ctx    context.Context
		amount int
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"加算できる",
			fields{balanceRepo, paymentRepo},
			args{ctx, 1, 10, 0},
			false,
		},
		{
			"減算できない",
			fields{balanceRepo, paymentRepo},
			args{ctx, 0, 10, 0},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				BalanceRepo: tt.fields.BalanceRepo,
				PaymentRepo: tt.fields.PaymentRepo,
			}
			if err := s.AddToUsers(tt.args.ctx, tt.args.amount, tt.args.limit, tt.args.offset); (err != nil) != tt.wantErr {
				t.Errorf("paymentService.AddToUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
