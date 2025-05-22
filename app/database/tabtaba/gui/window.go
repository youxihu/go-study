package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/google/wire"
)

var Set = wire.NewSet(NewWindow)

type Window struct {
	fyne.App
}

func NewWindow() *Window {
	myapp := app.New()
	return &Window{
		App: myapp,
	}
}
