package service

import (
	"context"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(AddService) AddService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type AddService interface {
	// [method=post,expose=true]
	Sum(ctx context.Context, a int64, b int64) (res int64, err error)
}

// the concrete implementation of service interface
type stubAddService struct {
	logger log.Logger `json:"logger"`
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(logger log.Logger) (s AddService) {
	var svc AddService
	{
		svc = &stubAddService{logger: logger}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Sum
func (ad *stubAddService) Sum(ctx context.Context, a int64, b int64) (res int64, err error) {
	return a + b, err
}
