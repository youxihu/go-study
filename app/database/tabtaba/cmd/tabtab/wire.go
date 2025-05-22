// cmd/wire.go
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-study/app/database/tabtaba/gui"
)

func InitApp() *gui.Window {
	wire.Build(
		gui.Set,
	)
	return nil
}
