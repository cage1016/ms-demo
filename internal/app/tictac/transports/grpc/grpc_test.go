package transports_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/cage1016/ms-sample/internal/app/tictac/endpoints"
	transports "github.com/cage1016/ms-sample/internal/app/tictac/transports/grpc"
	automocks "github.com/cage1016/ms-sample/internal/mocks/app/tictac/service"
	"github.com/cage1016/ms-sample/internal/pkg/errors"
	pb "github.com/cage1016/ms-sample/pb/tictac"
)

const (
	hostPort string = "localhost:8002"
)

func Test_Tic(t *testing.T) {
	type fields struct {
		svc *automocks.MockTictacService
	}
	type args struct {
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		wantErr   bool
		checkFunc func(err error)
	}{
		{
			name: "tic should return nil",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Tic(gomock.Any()).Return(nil),
				)
			},
			wantErr: false,
		},
		{
			name: "tic should return dummyError",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Tic(gomock.Any()).Return(errors.New("dummyError")),
				)
			},
			wantErr: true,
			checkFunc: func(err error) {
				s, _ := status.FromError(err)
				assert.Contains(t, s.Message(), "internal server error", fmt.Sprintf("tic should return dummyError: expected err %s contain %s", s.Message(), err.Error()))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				svc: automocks.NewMockTictacService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			logger := log.NewLogfmtLogger(os.Stderr)
			zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
			tracer := opentracing.GlobalTracer()

			// server
			server := grpc.NewServer()
			eps := endpoints.New(f.svc, logger, tracer, zkt)
			sc, err := net.Listen("tcp", hostPort)
			if err != nil {
				t.Fatalf("unable to listen: %+v", err)
			}
			defer server.GracefulStop()

			go func() {
				pb.RegisterTictacServer(server, transports.MakeGRPCServer(eps, tracer, zkt, logger))
				_ = server.Serve(sc)
			}()

			// client
			cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
			if err != nil {
				t.Fatalf("unable to Dial: %+v", err)
			}
			svc := transports.NewGRPCClient(cc, tracer, zkt, logger)

			if err := svc.Tic(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Tic(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(err)
				}
			}
		})
	}
}

func Test_Tac(t *testing.T) {
	type fields struct {
		svc *automocks.MockTictacService
	}
	type args struct {
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		wantErr   bool
		checkFunc func(res int64, err error)
	}{
		{
			name: "tac should return 5",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Tac(gomock.Any()).Return(int64(5), nil),
				)
			},
			wantErr: false,
			checkFunc: func(res int64, err error) {
				assert.Equal(t, int64(5), res, fmt.Sprintf("tac should return 5: expected 5 got %d", res))
			},
		},
		{
			name: "tic should return dummyError",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Tac(gomock.Any()).Return(int64(0), errors.New("dummyError")),
				)
			},
			wantErr: true,
			checkFunc: func(res int64, err error) {
				s, _ := status.FromError(err)
				assert.Contains(t, s.Message(), "internal server error", fmt.Sprintf("tic should return dummyError: expected err %s contain %s", s.Message(), err.Error()))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				svc: automocks.NewMockTictacService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			logger := log.NewLogfmtLogger(os.Stderr)
			zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
			tracer := opentracing.GlobalTracer()

			// server
			server := grpc.NewServer()
			eps := endpoints.New(f.svc, logger, tracer, zkt)
			sc, err := net.Listen("tcp", hostPort)
			if err != nil {
				t.Fatalf("unable to listen: %+v", err)
			}
			defer server.GracefulStop()

			go func() {
				pb.RegisterTictacServer(server, transports.MakeGRPCServer(eps, tracer, zkt, logger))
				_ = server.Serve(sc)
			}()

			// client
			cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
			if err != nil {
				t.Fatalf("unable to Dial: %+v", err)
			}
			svc := transports.NewGRPCClient(cc, tracer, zkt, logger)

			if res, err := svc.Tac(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Tac(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(res, err)
				}
			}
		})
	}
}
