package transports

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/rs/cors"

	"github.com/cage1016/ms-sample/internal/app/tictac/endpoints"
	"github.com/cage1016/ms-sample/internal/app/tictac/service"
	"github.com/cage1016/ms-sample/internal/pkg/errors"
	"github.com/cage1016/ms-sample/internal/pkg/jwt"
	"github.com/cage1016/ms-sample/internal/pkg/responses"
)

const (
	contentType string = "application/json"
)

// ShowTictac godoc
// @Summary Tic
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /tic [post]
func TicHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Post("/tic", httptransport.NewServer(
		endpoints.TicEndpoint,
		decodeHTTPTicRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Tic", logger), jwt.HTTPToContext()))...,
	))
}

// ShowTictac godoc
// @Summary Tac
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /tac [get]
func TacHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Get("/tac", httptransport.NewServer(
		endpoints.TacEndpoint,
		decodeHTTPTacRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Tac", logger), jwt.HTTPToContext()))...,
	))
}

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoints.Endpoints, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(responses.ErrorEncodeJSONResponse(CustomErrorEncoder)),
		httptransport.ServerErrorLogger(logger),
	}

	if zipkinTracer != nil {
		// Zipkin HTTP Server Trace can either be instantiated per endpoint with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit endpoint as ServerOption.
		// In the latter case, the operation name will be the endpoint's http method.
		// We demonstrate a global tracing service here.
		options = append(options, zipkin.HTTPServerTrace(zipkinTracer))
	}

	m := bone.New()
	TicHandler(m, endpoints, options, otTracer, logger)
	TacHandler(m, endpoints, options, otTracer, logger)
	return cors.AllowAll().Handler(m)
}

// decodeHTTPTicRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPTicRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.TicRequest
	return req, nil
}

// decodeHTTPTacRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPTacRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.TacRequest
	return req, nil
}

// NewHTTPClient returns an AddService backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (service.TictacService, error) { // Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// global client middlewares
	var options []httptransport.ClientOption

	if zipkinTracer != nil {
		// Zipkin HTTP Client Trace can either be instantiated per endpoint with a
		// provided operation name or a global tracing client can be instantiated
		// without an operation name and fed to each Go kit endpoint as ClientOption.
		// In the latter case, the operation name will be the endpoint's http method.
		options = append(options, zipkin.HTTPClientTrace(zipkinTracer))
	}

	e := endpoints.Endpoints{}

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	// The Tic endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var ticEndpoint endpoint.Endpoint
	{
		ticEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/tic"),
			encodeHTTPTicRequest,
			decodeHTTPTicResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		ticEndpoint = opentracing.TraceClient(otTracer, "Tic")(ticEndpoint)
		ticEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Tic")(ticEndpoint)
		e.TicEndpoint = ticEndpoint
	}

	// The Tac endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var tacEndpoint endpoint.Endpoint
	{
		tacEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/tac"),
			encodeHTTPTacRequest,
			decodeHTTPTacResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		tacEndpoint = opentracing.TraceClient(otTracer, "Tac")(tacEndpoint)
		tacEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Tac")(tacEndpoint)
		e.TacEndpoint = tacEndpoint
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return e, nil
}

//
func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// encodeHTTPTicRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPTicRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPTicResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a client.
func decodeHTTPTicResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, responses.JSONErrorDecoder(r)
	}
	var resp endpoints.TicResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// encodeHTTPTacRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPTacRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPTacResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a client.
func decodeHTTPTacResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, responses.JSONErrorDecoder(r)
	}
	var resp endpoints.TacResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func CustomErrorEncoder(errorVal errors.Error) (code int) {
	switch {
	// TODO write your own custom error check here
	case errors.Contains(errorVal, jwt.ErrXJWTContextMissing):
		code = http.StatusForbidden
	}
	return
}
