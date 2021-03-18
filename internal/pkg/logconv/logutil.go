package logconv

import (
	"github.com/go-kit/kit/log/level"
)

const (
	Debug = "debug"
	Info  = "info"
	Warn  = "warn"
	Error = "error"
)

var ls = map[string]level.Option{
	Debug: level.AllowDebug(),
	Info:  level.AllowInfo(),
	Warn:  level.AllowWarn(),
	Error: level.AllowError(),
}

func Atol(l string) level.Option {
	return ls[l]
}
