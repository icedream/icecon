// This file was created by ui2walk and may be regenerated.
// DO NOT EDIT OR YOUR MODIFICATIONS WILL BE LOST!

package main

import (
	"github.com/lxn/walk"
)

type connectDialogUI struct {
	label_2      *walk.Label
	rconAddress  *walk.LineEdit
	rconPassword *walk.LineEdit
	label        *walk.Label
	cancel       *walk.PushButton
	ok           *walk.PushButton
}

func (w *connectDialog) init(owner walk.Form) (err error) {
	if w.Dialog, err = walk.NewDialog(owner); err != nil {
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

	w.SetName("connectDialog")
	if err := w.SetClientSize(walk.Size{400, 90}); err != nil {
		return err
	}
	if err := w.SetTitle(`Connect to...`); err != nil {
		return err
	}

	// label_2
	if w.ui.label_2, err = walk.NewLabel(w); err != nil {
		return err
	}
	w.ui.label_2.SetName("label_2")
	if err := w.ui.label_2.SetBounds(walk.Rectangle{10, 36, 50, 16}); err != nil {
		return err
	}
	if err := w.ui.label_2.SetText(`Password:`); err != nil {
		return err
	}

	// rconAddress
	if w.ui.rconAddress, err = walk.NewLineEdit(w); err != nil {
		return err
	}
	w.ui.rconAddress.SetName("rconAddress")
	if err := w.ui.rconAddress.SetBounds(walk.Rectangle{104, 10, 291, 20}); err != nil {
		return err
	}

	// rconPassword
	if w.ui.rconPassword, err = walk.NewLineEdit(w); err != nil {
		return err
	}
	w.ui.rconPassword.SetName("rconPassword")
	if err := w.ui.rconPassword.SetBounds(walk.Rectangle{104, 36, 291, 20}); err != nil {
		return err
	}
	w.ui.rconPassword.SetPasswordMode(true)

	// label
	if w.ui.label, err = walk.NewLabel(w); err != nil {
		return err
	}
	w.ui.label.SetName("label")
	if err := w.ui.label.SetBounds(walk.Rectangle{10, 10, 88, 16}); err != nil {
		return err
	}
	if err := w.ui.label.SetText(`Address (IP:Port):`); err != nil {
		return err
	}

	// cancel
	if w.ui.cancel, err = walk.NewPushButton(w); err != nil {
		return err
	}
	w.ui.cancel.SetName("cancel")
	if err := w.ui.cancel.SetBounds(walk.Rectangle{239, 60, 75, 23}); err != nil {
		return err
	}
	if err := w.ui.cancel.SetText(`Cancel`); err != nil {
		return err
	}

	// ok
	if w.ui.ok, err = walk.NewPushButton(w); err != nil {
		return err
	}
	w.ui.ok.SetName("ok")
	if err := w.ui.ok.SetBounds(walk.Rectangle{320, 60, 75, 23}); err != nil {
		return err
	}
	if err := w.ui.ok.SetText(`OK`); err != nil {
		return err
	}

	// Tab order

	succeeded = true

	return nil
}
