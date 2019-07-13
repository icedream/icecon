//+build windows

package main

import "github.com/lxn/walk"

type connectDialog struct {
	*walk.Dialog
	ui connectDialogUI
}

func runConnectDialog(currentAddr string, currentPw string, owner walk.Form, title string) (result bool, addr string, pw string, err error) {
	dlg := new(connectDialog)

	if err = dlg.init(owner); err != nil {
		return
	}

	if err = dlg.SetDefaultButton(dlg.ui.ok); err != nil {
		return
	}
	dlg.ui.ok.Clicked().Attach(func() {
		addr = dlg.ui.rconAddress.Text()
		pw = dlg.ui.rconPassword.Text()
		dlg.Accept()
	})

	if err = dlg.SetCancelButton(dlg.ui.cancel); err != nil {
		return
	}
	dlg.ui.cancel.Clicked().Attach(func() {
		dlg.Cancel()
	})

	dlg.ui.rconAddress.SetText(currentAddr)
	dlg.ui.rconPassword.SetText(currentPw)

	if title != "" {
		dlg.SetTitle(title)
	}

	choice := dlg.Run()

	result = choice == walk.DlgCmdOK

	return
}
