package transports

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cage1016/ms-sample/internal/app/tictac/endpoints"
	"github.com/cage1016/ms-sample/internal/app/tictac/service"
	"github.com/cage1016/ms-sample/internal/pkg/errors"
	"github.com/cage1016/ms-sample/internal/pkg/jwt"
	pb "github.com/cage1016/ms-sample/pb/tictac"
)

type grpcServer struct {
	tic grpctransport.Handler `json:""`
	tac grpctransport.Handler `json:""`
}

func (s *grpcServer) Tic(ctx context.Context, req *pb.TicRequest) (rep *pb.TicResponse, err error) {
	_, rp, err := s.tic.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.TicResponse)
	return rep, nil
}

func (s *grpcServer) Tac(ctx context.Context, req *pb.TacRequest) (rep *pb.TacResponse, err error) {
	_, rp, err := s.tac.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.TacResponse)
	return rep, nil
}

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints endpoints.Endpoints, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (req pb.TictacServer) { // Zipkin GRPC Server Trace can either be instantiated per gRPC method with a
	// provided operation name or a global tracing service can be instantiated
	// without an operation name and fed to each Go kit gRPC server as a
	// ServerOption.
	// In the latter case, the operation name will be the endpoint's grpc method
	// path if used in combination with the Go kit gRPC Interceptor.
	//
	// In this example, we demonstrate a global Zipkin tracing service with
	// Go kit gRPC Interceptor.
	zipkinServer := zipkin.GRPCServerTrace(zipkinTracer)

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
		zipkinServer,
	}

	return &grpcServer{
		tic: grpctransport.NewServer(
			endpoints.TicEndpoint,
			decodeGRPCTicRequest,
			encodeGRPCTicResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Tic", logger), jwt.GRPCToContext()))...,
		),

		tac: grpctransport.NewServer(
			endpoints.TacEndpoint,
			decodeGRPCTacRequest,
			encodeGRPCTacResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Tac", logger), jwt.GRPCToContext()))...,
		),
	}
}

// decodeGRPCTicRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCTicRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	_ = grpcReq.(*pb.TicRequest)
	return endpoints.TicRequest{}, nil
}

// encodeGRPCTicResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCTicResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.TicResponse)
	return &pb.TicResponse{}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCTacRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCTacRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	_ = grpcReq.(*pb.TacRequest)
	return endpoints.TacRequest{}, nil
}

// encodeGRPCTacResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCTacResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.TacResponse)
	return &pb.TacResponse{Res: reply.Res}, grpcEncodeError(errors.Cast(reply.Err))
}

// NewGRPCClient returns an AddService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.TictacService { // Zipkin GRPC Client Trace can either be instantiated per gRPC method with a
	// provided operation name or a global tracing client can be instantiated
	// without an operation name and fed to each Go kit client as ClientOption.
	// In the latter case, the operation name will be the endpoint's grpc method
	// path.
	//
	// In this example, we demonstrace a global tracing client.
	zipkinClient := zipkin.GRPCClientTrace(zipkinTracer)

	// global client middlewares
	options := []grpctransport.ClientOption{
		zipkinClient,
	}

	// The Tic endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var ticEndpoint endpoint.Endpoint
	{
		ticEndpoint = grpctransport.NewClient(
			conn,
			"pb.Tictac",
			"Tic",
			encodeGRPCTicRequest,
			decodeGRPCTicResponse,
			pb.TicResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), jwt.ContextToGRPC()))...,
		).Endpoint()
		ticEndpoint = opentracing.TraceClient(otTracer, "Tic")(ticEndpoint)
	}

	// The Tac endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var tacEndpoint endpoint.Endpoint
	{
		tacEndpoint = grpctransport.NewClient(
			conn,
			"pb.Tictac",
			"Tac",
			encodeGRPCTacRequest,
			decodeGRPCTacResponse,
			pb.TacResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), jwt.ContextToGRPC()))...,
		).Endpoint()
		tacEndpoint = opentracing.TraceClient(otTracer, "Tac")(tacEndpoint)
	}

	return endpoints.Endpoints{
		TicEndpoint: ticEndpoint,
		TacEndpoint: tacEndpoint,
	}
}

// encodeGRPCTicRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Tic request to a gRPC Tic request. Primarily useful in a client.
func encodeGRPCTicRequest(_ context.Context, request interface{}) (interface{}, error) {
	_ = request.(endpoints.TicRequest)
	return &pb.TicRequest{}, nil
}

// decodeGRPCTicResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Tic reply to a user-domain Tic response. Primarily useful in a client.
func decodeGRPCTicResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	_ = grpcReply.(*pb.TicResponse)
	return endpoints.TicResponse{}, nil
}

// encodeGRPCTacRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Tac request to a gRPC Tac request. Primarily useful in a client.
func encodeGRPCTacRequest(_ context.Context, request interface{}) (interface{}, error) {
	_ = request.(endpoints.TacRequest)
	return &pb.TacRequest{}, nil
}

// decodeGRPCTacResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Tac reply to a user-domain Tac response. Primarily useful in a client.
func decodeGRPCTacResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.TacResponse)
	return endpoints.TacResponse{Res: reply.Res}, nil
}

func grpcEncodeError(err errors.Error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		return status.Error(st.Code(), st.Message())
	}

	switch {
	// TODO write your own custom error check here
	case errors.Contains(err, jwt.ErrXJWTContextMissing):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
