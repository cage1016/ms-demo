package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	addservice "github.com/cage1016/ms-sample/internal/app/add/service"
	"github.com/cage1016/ms-sample/internal/app/tictac/model"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(TictacService) TictacService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
//go:generate mockgen -destination ../../../../internal/mocks/app/tictac/service/tictacservice.go -package=automocks . TictacService
type TictacService interface {
	// [method=post,expose=true]
	Tic(ctx context.Context) (err error)
	// [method=get,expose=true]
	Tac(ctx context.Context) (res int64, err error)
}

// the concrete implementation of service interface
type stubTictacService struct {
	logger log.Logger
	addsvc addservice.AddService
	repo   model.TictacRespository
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(repo model.TictacRespository, addsvc addservice.AddService, logger log.Logger) (s TictacService) {
	var svc TictacService
	{
		svc = &stubTictacService{logger: logger, addsvc: addsvc, repo: repo}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Tic
func (ti *stubTictacService) Tic(ctx context.Context) (err error) {
	dv, err := ti.repo.Tac(ctx)
	if err != nil {
		level.Error(ti.logger).Log("method", "ti.repo.Tac", "err", err)
		return err
	}

	nv, err := ti.addsvc.Sum(ctx, dv, 1)
	if err != nil {
		level.Error(ti.logger).Log("method", "ti.addsvc.Sum", "err", err)
		return err
	}

	return ti.repo.Tic(ctx, nv)
}

// Implement the business logic of Tac
func (ti *stubTictacService) Tac(ctx context.Context) (res int64, err error) {
	return ti.repo.Tac(ctx)
}
