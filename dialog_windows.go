//+build windows

package main

import "github.com/lxn/walk"

type mainDialog struct {
	*walk.MainWindow
	ui mainDialogUI
}
