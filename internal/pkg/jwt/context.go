package jwt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc/metadata"

	"github.com/cage1016/ms-sample/internal/pkg/errors"
)

var (
	ErrXJwtPlayloadContextMissing = errors.New(fmt.Sprintf("%s not passed through the context", XJwtPlayload))
)

const XJwtPlayload = "X-Jwt-Playload"

func HTTPToContext() kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		xjp := r.Header.Get(XJwtPlayload)
		if len(xjp) == 0 {
			return ctx
		}
		return context.WithValue(ctx, XJwtPlayload, xjp)
	}
}

func ContextToHTTP() kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		version, ok := ctx.Value(XJwtPlayload).(string)
		if ok {
			r.Header.Add(XJwtPlayload, version)
		}
		return ctx
	}
}

func GRPCToContext() grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		// capital "Key" is illegal in HTTP/2.
		xJwtPlayloadHeader, ok := md["x-jwt-playload"]
		if !ok {
			return ctx
		}

		xjp := xJwtPlayloadHeader[0]
		if len(xjp) == 0 {
			return ctx
		}

		fmt.Println("GRPCToContext: xjp", xjp)
		return context.WithValue(ctx, XJwtPlayload, xjp)
	}
}

func ContextToGRPC() grpc.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		xjp, ok := ctx.Value(XJwtPlayload).(string)
		if ok {
			fmt.Println("ContextToGRPC: xjp", xjp)
			// capital "Key" is illegal in HTTP/2.
			(*md)["x-jwt-playload"] = []string{xjp}
		}

		return ctx
	}
}
