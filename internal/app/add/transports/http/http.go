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

	"github.com/cage1016/ms-sample/internal/app/add/endpoints"
	"github.com/cage1016/ms-sample/internal/app/add/service"
	"github.com/cage1016/ms-sample/internal/pkg/errors"
	"github.com/cage1016/ms-sample/internal/pkg/jwt"
	"github.com/cage1016/ms-sample/internal/pkg/responses"
)

// ShowAdd godoc
// @Summary Sum
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /sum [post]
func SumHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Post("/sum", httptransport.NewServer(
		endpoints.SumEndpoint,
		decodeHTTPSumRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Sum", logger), jwt.HTTPToContext()))...,
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
	SumHandler(m, endpoints, options, otTracer, logger)
	return cors.AllowAll().Handler(m)
}

// decodeHTTPSumRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.SumRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// NewHTTPClient returns an AddService backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (service.AddService, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// global client middlewares
	options := []httptransport.ClientOption{}

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
	// The Sum endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/sum"),
			encodeHTTPSumRequest,
			decodeHTTPSumResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		sumEndpoint = opentracing.TraceClient(otTracer, "Sum")(sumEndpoint)
		if zipkinTracer != nil {
			sumEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Sum")(sumEndpoint)
		}
		e.SumEndpoint = sumEndpoint
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

// encodeHTTPSumRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPSumRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPSumResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a client.
func decodeHTTPSumResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, responses.JSONErrorDecoder(r)
	}
	var resp endpoints.SumResponse
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
