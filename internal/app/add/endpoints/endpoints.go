package endpoints

import (
	"context"

	"github.com/cage1016/ms-sample/internal/app/add/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
)

// Endpoints collects all of the endpoints that compose the add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	SumEndpoint endpoint.Endpoint `json:""`
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.AddService, logger log.Logger, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) (ep Endpoints) {
	var sumEndpoint endpoint.Endpoint
	{
		method := "sum"
		sumEndpoint = MakeSumEndpoint(svc)
		sumEndpoint = opentracing.TraceServer(otTracer, method)(sumEndpoint)
		sumEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(sumEndpoint)
		sumEndpoint = LoggingMiddleware(log.With(logger, "method", method))(sumEndpoint)
		ep.SumEndpoint = sumEndpoint
	}

	return ep
}

// MakeSumEndpoint returns an endpoint that invokes Sum on the service.
// Primarily useful in a server.
func MakeSumEndpoint(svc service.AddService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SumRequest)
		if err := req.validate(); err != nil {
			return SumResponse{}, err
		}
		res, err := svc.Sum(ctx, req.A, req.B)
		return SumResponse{Res: res}, err
	}
}

// Sum implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Sum(ctx context.Context, a int64, b int64) (res int64, err error) {
	resp, err := e.SumEndpoint(ctx, SumRequest{A: a, B: b})
	if err != nil {
		return
	}
	response := resp.(SumResponse)
	return response.Res, nil
}
