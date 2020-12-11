package transports_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/stretchr/testify/assert"

	"github.com/cage1016/ms-sample/internal/app/tictac/endpoints"
	transports "github.com/cage1016/ms-sample/internal/app/tictac/transports/http"
	automocks "github.com/cage1016/ms-sample/internal/mocks/app/tictac/service"
	test "github.com/cage1016/ms-sample/test/util"
)

func Test_Toc(t *testing.T) {
	type fields struct {
		svc *automocks.MockTictacService
	}
	type args struct {
		method, url string
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res *http.Response, err error, body []byte)
	}{
		{
			name: "tic",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Tic(gomock.Any()).Return(nil),
				)
			},
			wantErr: false,
			args: args{
				method: http.MethodPost,
				url:    "/tic",
			},
			checkFunc: func(res *http.Response, err error, body []byte) {
				assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
				assert.Equal(t, res.StatusCode, 204, fmt.Sprintf("status should be 204: got %d", res.StatusCode))
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

			eps := endpoints.New(f.svc, logger, tracer, zkt)
			ts := httptest.NewServer(transports.NewHTTPHandler(eps, tracer, zkt, logger))
			defer ts.Close()

			req := test.TestRequest{
				Client:      ts.Client(),
				Method:      tt.args.method,
				URL:         fmt.Sprintf("%s%s", ts.URL, tt.args.url),
				ContentType: "application/json",
			}

			if res, err := req.Make(); (err != nil) != tt.wantErr {
				t.Errorf("%s: unexpected error %s", tt.name, err)
			} else {
				body, _ := ioutil.ReadAll(res.Body)
				if tt.checkFunc != nil {
					tt.checkFunc(res, err, body)
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
		method, url string
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res *http.Response, err error, body []byte)
	}{
		{
			name: "tac should contain 2",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Tac(gomock.Any()).Return(int64(2), nil),
				)
			},
			wantErr: false,
			args: args{
				method: http.MethodGet,
				url:    "/tac",
			},
			checkFunc: func(res *http.Response, err error, body []byte) {
				assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
				assert.Contains(t, string(body), "2")
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

			eps := endpoints.New(f.svc, logger, tracer, zkt)
			ts := httptest.NewServer(transports.NewHTTPHandler(eps, tracer, zkt, logger))
			defer ts.Close()

			req := test.TestRequest{
				Client:      ts.Client(),
				Method:      tt.args.method,
				URL:         fmt.Sprintf("%s%s", ts.URL, tt.args.url),
				ContentType: "application/json",
			}

			if res, err := req.Make(); (err != nil) != tt.wantErr {
				t.Errorf("%s: unexpected error %s", tt.name, err)
			} else {
				body, _ := ioutil.ReadAll(res.Body)
				if tt.checkFunc != nil {
					tt.checkFunc(res, err, body)
				}
			}
		})
	}
}
