package model

import "context"

type TicTac struct {
	Value int64 `json:"value" db:"value"`
}

//go:generate mockgen -destination ../../../../internal/mocks/app/tictac/model/tictacrespository.go -package=automocks . TicTacRespository
type TicTacRespository interface {
	Add(context.Context) error
	Get(context.Context) (int64, error)
}
