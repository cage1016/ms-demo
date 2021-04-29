package model

import "context"

type Tictac struct {
	Value int64
}

//go:generate mockgen -destination ../../../../internal/mocks/app/tictac/model/tictacrespository.go -package=automocks . TictacRespository
type TictacRespository interface {
	Tic(context.Context, int64) error
	Tac(context.Context) (int64, error)
}
