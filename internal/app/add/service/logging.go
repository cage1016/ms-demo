package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingMiddleware struct {
	logger log.Logger `json:""`
	next   AddService `json:""`
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next AddService) AddService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (lm loggingMiddleware) Sum(ctx context.Context, a int64, b int64) (res int64, err error) {
	defer func() {
		lm.logger.Log("method", "Sum", "a", a, "b", b, "err", err)
	}()

	return lm.next.Sum(ctx, a, b)
}
