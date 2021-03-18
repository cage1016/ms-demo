// +build !integration

package postgres_test

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/kit/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	psql "github.com/cage1016/ms-sample/internal/app/tictac/postgres"
)

func Test_Tic(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		wantErr bool
	}{
		{
			name: "Tic func",
			prepare: func(f *fields) {
				f.mock.ExpectExec(regexp.QuoteMeta("UPDATE \"tictacs\" SET \"value\"=value+1 WHERE 1=1")).WillReturnResult(sqlmock.NewResult(1, 1))
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

			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := psql.New(gdb, log.NewLogfmtLogger(os.Stderr))

			if err := repo.Tic(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Add(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Tac(t *testing.T) {
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
				f.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"tictacs\"")).WillReturnRows(rows).WillReturnError(nil)
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

			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := psql.New(gdb, log.NewLogfmtLogger(os.Stderr))

			if res, err := repo.Tac(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Get(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log(res)
			}
		})
	}
}
