package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingMiddleware struct {
	logger log.Logger    `json:""`
	next   TictacService `json:""`
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next TictacService) TictacService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (lm loggingMiddleware) Tic(ctx context.Context) (err error) {
	defer func() {
		lm.logger.Log("method", "Tic", "err", err)
	}()

	return lm.next.Tic(ctx)
}

func (lm loggingMiddleware) Tac(ctx context.Context) (res int64, err error) {
	defer func() {
		lm.logger.Log("method", "Tac", "err", err)
	}()

	return lm.next.Tac(ctx)
}
