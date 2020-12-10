package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"

	"github.com/cage1016/ms-sample/internal/app/tictac/postgres"
)

func Test_Add(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		wantErr bool
	}{
		{
			name: "Add func",
			prepare: func(f *fields) {
				f.mock.ExpectExec("update tictac").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			f := fields{
				mock: mock,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			repo := postgres.New(sqlx.NewDb(db, "postgres"), log.NewLogfmtLogger(os.Stderr))

			if err := repo.Add(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Add(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Get(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		wantErr bool
	}{
		{
			name: "Get func",
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"value"}).AddRow(2)
				f.mock.ExpectQuery("^select (.+) from tictac").WillReturnRows(rows)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			f := fields{
				mock: mock,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			repo := postgres.New(sqlx.NewDb(db, "postgres"), log.NewLogfmtLogger(os.Stderr))

			if res, err := repo.Get(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Get(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log(res)
			}
		})
	}
}
