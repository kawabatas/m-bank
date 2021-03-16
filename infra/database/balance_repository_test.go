package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kawabatas/m-bank/domain/model"
)

func newBalanceRepo(t *testing.T) *BalanceRepository {
	t.Helper()
	db := newTestConnection(t)
	return NewBalanceRepository(db)
}

func TestBalanceRepository_Get(t *testing.T) {
	repo := newBalanceRepo(t)
	users := createSampleUsers(t, repo.DB, 1)
	ctx := context.Background()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Balance
		wantErr bool
	}{
		{
			"取得できる",
			fields{repo.DB},
			args{ctx, users[0].ID},
			&model.Balance{
				UserID: users[0].ID,
				Amount: initBalanceAmount,
			},
			false,
		},
		{
			"取得できない",
			fields{repo.DB},
			args{ctx, users[0].ID + 100},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BalanceRepository{
				DB: tt.fields.DB,
			}
			got, err := r.Get(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, nil); diff != "" {
				t.Errorf("BalanceRepository.Get() mismatch (-want +got): \n %s", diff)
			}
		})
	}
}

func TestBalanceRepository_AddToUsers(t *testing.T) {
	repo := newBalanceRepo(t)
	users := createSampleUsers(t, repo.DB, 3)
	ctx := context.Background()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx    context.Context
		amount int
		limit  int
		offset int
	}
	type want struct {
		Amounts   []int
		LogCounts []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			"全員の残高へ+1",
			fields{repo.DB},
			args{ctx, 1, 10, 0},
			want{
				[]int{initBalanceAmount + 1, initBalanceAmount + 1, initBalanceAmount + 1},
				[]int{1, 1, 1},
			},
			false,
		},
		{
			"limit,offsetを指定して、残高へ+10",
			fields{repo.DB},
			args{ctx, 10, 1, 1},
			want{
				[]int{initBalanceAmount + 1, initBalanceAmount + 1 + 10, initBalanceAmount + 1},
				[]int{1, 2, 1},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BalanceRepository{
				DB: tt.fields.DB,
			}
			if err := r.AddToUsers(tt.args.ctx, tt.args.amount, tt.args.limit, tt.args.offset); (err != nil) != tt.wantErr {
				t.Errorf("BalanceRepository.AddToUsers() error = %v, wantErr %v", err, tt.wantErr)
			}

			var amounts, logCounts []int
			for _, u := range users {
				b, err := r.Get(context.Background(), u.ID)
				if err != nil {
					t.Errorf("BalanceRepository.AddToUsers() r.Get error %v", err)
				}
				amounts = append(amounts, int(b.Amount))
				c, err := countBalanceLog(context.Background(), r.DB, u.ID)
				if err != nil {
					t.Errorf("BalanceRepository.AddToUsers() countBalanceLog error %v", err)
				}
				logCounts = append(logCounts, c)
			}
			got := want{Amounts: amounts, LogCounts: logCounts}

			if diff := cmp.Diff(tt.want, got, nil); diff != "" {
				t.Errorf("BalanceRepository.AddToUsers() mismatch (-want +got): \n %s", diff)
			}
		})
	}
}
