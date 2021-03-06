// +build !integration

package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"

	"github.com/cage1016/ms-sample/internal/app/tictac/service"
	automocks2 "github.com/cage1016/ms-sample/internal/mocks/app/add/service"
	automocks "github.com/cage1016/ms-sample/internal/mocks/app/tictac/model"
)

func Test_Tic(t *testing.T) {
	type fields struct {
		repo   *automocks.MockTictacRespository
		addsvc *automocks2.MockAddService
	}
	type args struct {
		value int64
	}

	tests := []struct {
		name      string
		args      args
		prepare   func(f *fields)
		wantErr   bool
		checkFunc func()
	}{
		{
			name: "tic should return nil",
			args: args{1},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Tac(context.Background()).Return(int64(1), nil),
					f.addsvc.EXPECT().Sum(context.Background(), int64(1), int64(1)).Return(int64(2), nil),
					f.repo.EXPECT().Tic(context.Background(), int64(2)).Return(nil),
				)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:   automocks.NewMockTictacRespository(ctrl),
				addsvc: automocks2.NewMockAddService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.New(f.repo, f.addsvc, log.NewLogfmtLogger(os.Stderr))

			if err := s.Tic(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Tic(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc()
				}
			}
		})
	}
}

func Test_Tac(t *testing.T) {
	type fields struct {
		repo   *automocks.MockTictacRespository
		addsvc *automocks2.MockAddService
	}
	type args struct {
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res int64, err error)
	}{
		{
			name: "tac should return 2",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Tac(context.Background()).Return(int64(2), nil),
				)
			},
			wantErr: false,
			checkFunc: func(res int64, err error) {
				assert.Equal(t, int64(2), res, fmt.Sprintf("tac should return 2: expected 2 got %d", res))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:   automocks.NewMockTictacRespository(ctrl),
				addsvc: automocks2.NewMockAddService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := service.New(f.repo, f.addsvc, log.NewLogfmtLogger(os.Stderr))

			if res, err := s.Tac(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Tic(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(res, err)
				}
			}
		})
	}
}
