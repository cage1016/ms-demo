package model

import "context"

type TicTac struct {
	Value int64 `json:"value" db:"value"`
}

type TicTacRespository interface {
	Add(context.Context) error
	Get(context.Context) (int64, error)
}
