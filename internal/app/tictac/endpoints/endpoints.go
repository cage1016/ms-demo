package endpoints

import (
	"context"

	"github.com/cage1016/ms-sample/internal/app/tictac/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
)

// Endpoints collects all of the endpoints that compose the tictac service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	TicEndpoint endpoint.Endpoint `json:""`
	TacEndpoint endpoint.Endpoint `json:""`
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.TictacService, logger log.Logger, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) (ep Endpoints) {
	var ticEndpoint endpoint.Endpoint
	{
		method := "tic"
		ticEndpoint = MakeTicEndpoint(svc)
		ticEndpoint = opentracing.TraceServer(otTracer, method)(ticEndpoint)
		if zipkinTracer != nil {
			ticEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(ticEndpoint)
		}
		ticEndpoint = LoggingMiddleware(log.With(logger, "method", method))(ticEndpoint)
		ep.TicEndpoint = ticEndpoint
	}

	var tacEndpoint endpoint.Endpoint
	{
		method := "tac"
		tacEndpoint = MakeTacEndpoint(svc)
		tacEndpoint = opentracing.TraceServer(otTracer, method)(tacEndpoint)
		if zipkinTracer != nil {
			tacEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(tacEndpoint)
		}
		tacEndpoint = LoggingMiddleware(log.With(logger, "method", method))(tacEndpoint)
		ep.TacEndpoint = tacEndpoint
	}

	return ep
}

// MakeTicEndpoint returns an endpoint that invokes Tic on the service.
// Primarily useful in a server.
func MakeTicEndpoint(svc service.TictacService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		err := svc.Tic(ctx)
		return TicResponse{}, err
	}
}

// Tic implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Tic(ctx context.Context) (err error) {
	resp, err := e.TicEndpoint(ctx, TicRequest{})
	if err != nil {
		return
	}
	_ = resp.(TicResponse)
	return nil
}

// MakeTacEndpoint returns an endpoint that invokes Tac on the service.
// Primarily useful in a server.
func MakeTacEndpoint(svc service.TictacService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		res, err := svc.Tac(ctx)
		return TacResponse{Res: res}, err
	}
}

// Tac implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Tac(ctx context.Context) (res int64, err error) {
	resp, err := e.TacEndpoint(ctx, TacRequest{})
	if err != nil {
		return
	}
	response := resp.(TacResponse)
	return response.Res, nil
}
