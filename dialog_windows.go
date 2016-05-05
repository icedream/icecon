//+build windows
//go:generate ui2walk dialog.ui

package main

import (
	"github.com/lxn/walk"
)

type mainDialog struct {
	*walk.MainWindow
	ui mainDialogUI
}
