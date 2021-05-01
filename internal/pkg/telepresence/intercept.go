package telepresence

import (
	"context"
	"fmt"
	stdhttp "net/http"

	"github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc/metadata"
)

const XTelepresenceInterceptId = "x-telepresence-intercept-id"

func HTTPToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		interceptId := r.Header.Get(XTelepresenceInterceptId)
		if len(interceptId) == 0 {
			return ctx
		}
		fmt.Println("HTTPToContext", interceptId)
		return context.WithValue(ctx, XTelepresenceInterceptId, interceptId)
	}
}

func ContextToHTTP() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		interceptId, ok := ctx.Value(XTelepresenceInterceptId).(string)
		if ok {
			r.Header.Add(XTelepresenceInterceptId, interceptId)
		}
		fmt.Println("ContextToHTTP", interceptId)
		return ctx
	}
}

func GRPCToContext() grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		// capital "Key" is illegal in HTTP/2.
		playloadHeader, ok := md["x-telepresence-intercept-id"]
		if !ok {
			return ctx
		}

		interceptId := playloadHeader[0]
		if len(interceptId) == 0 {
			return ctx
		}
		fmt.Println("GRPCToContext", interceptId)
		return context.WithValue(ctx, XTelepresenceInterceptId, interceptId)
	}
}

func ContextToGRPC() grpc.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		interceptId, ok := ctx.Value(XTelepresenceInterceptId).(string)
		if ok {
			fmt.Println("ContextToGRPC", interceptId)
			// capital "Key" is illegal in HTTP/2.
			(*md)["x-telepresence-intercept-id"] = []string{interceptId}
		}

		return ctx
	}
}
