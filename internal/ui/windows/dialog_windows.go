//go:build windows
// +build windows

package windows

import "github.com/lxn/walk"

type mainDialog struct {
	*walk.MainWindow
	ui mainDialogUI
}
