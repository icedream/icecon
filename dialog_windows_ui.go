// This file was created by ui2walk and may be regenerated.
// DO NOT EDIT OR YOUR MODIFICATIONS WILL BE LOST!

package main

import (
	"github.com/lxn/walk"
)

type mainDialogUI struct {
	centralwidget *walk.Composite
	rconOutput    *walk.TextEdit
	rconInput     *walk.LineEdit
}

func (w *mainDialog) init() (err error) {
	if w.MainWindow, err = walk.NewMainWindow(); err != nil {
		return err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	var font *walk.Font
	if font == nil {
		font = nil
	}

	w.SetName("mainDialog")
	l := walk.NewVBoxLayout()
	if err := l.SetMargins(walk.Margins{0, 0, 0, 0}); err != nil {
		return err
	}
	if err := w.SetLayout(l); err != nil {
		return err
	}
	if err := w.SetClientSize(walk.Size{800, 560}); err != nil {
		return err
	}
	if err := w.SetTitle(`IceCon`); err != nil {
		return err
	}

	// centralwidget
	if w.ui.centralwidget, err = walk.NewComposite(w); err != nil {
		return err
	}
	w.ui.centralwidget.SetName("centralwidget")
	verticalLayout := walk.NewVBoxLayout()
	if err := w.ui.centralwidget.SetLayout(verticalLayout); err != nil {
		return err
	}
	if err := verticalLayout.SetMargins(walk.Margins{9, 9, 9, 9}); err != nil {
		return err
	}
	if err := verticalLayout.SetSpacing(6); err != nil {
		return err
	}

	// rconOutput
	if w.ui.rconOutput, err = walk.NewTextEdit(w.ui.centralwidget); err != nil {
		return err
	}
	w.ui.rconOutput.SetName("rconOutput")
	if font, err = walk.NewFont("Courier New", 8, 0); err != nil {
		return err
	}
	w.ui.rconOutput.SetFont(font)
	w.ui.rconOutput.SetReadOnly(true)

	// rconInput
	if w.ui.rconInput, err = walk.NewLineEdit(w.ui.centralwidget); err != nil {
		return err
	}
	w.ui.rconInput.SetName("rconInput")
	if font, err = walk.NewFont("Courier New", 8, 0); err != nil {
		return err
	}
	w.ui.rconInput.SetFont(font)
	if err := w.ui.rconInput.SetMinMaxSize(walk.Size{0, 0}, walk.Size{16777215, 21}); err != nil {
		return err
	}

	// Tab order

	succeeded = true

	return nil
}
