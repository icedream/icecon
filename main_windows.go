//+build windows

package main

//go:generate ui2walk

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/lxn/walk"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	hasGraphicalUI = true

	flagGui = kingpin.
		Flag("gui", "Run as GUI (runs automatically as GUI if no arguments given, ignored if command flag used)").
		Short('g').Bool()

	guiInitErr  error
	kernel32    *syscall.DLL
	freeConsole *syscall.Proc

	dlg *mainDialog
)

func init() {
	kernel32, guiInitErr = syscall.LoadDLL("kernel32.dll")
	freeConsole, guiInitErr = kernel32.FindProc("FreeConsole")
}

func uiLogError(text string) {
	uiNormalize(&text)
	dlg.Synchronize(func() {
		dlg.ui.rconOutput.AppendText("ERROR: " + text + "\r\n")
		walk.MsgBox(dlg, "Error",
			text,
			walk.MsgBoxIconError)
	})
}

func uiLog(text string) {
	uiNormalize(&text)
	dlg.Synchronize(func() {
		dlg.ui.rconOutput.AppendText(text + "\r\n")
	})
}

func uiNormalize(textRef *string) {
	text := *textRef

	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n", "\r\n", -1)

	*textRef = text
}

func runGraphicalUi() (err error) {
	dlg = new(mainDialog)
	if err := dlg.init(); err != nil {
		panic(err)
	}
	defer dlg.Dispose()

	// Window icon
	// TODO - Do this more intelligently
	for i := uintptr(0); i < uintptr(128); i++ {
		if icon, err := walk.NewIconFromResourceId(i); err == nil {
			dlg.SetIcon(icon)
			break
		}
	}

	// Quit button
	quitAction := walk.NewAction()
	if err = quitAction.SetText("&Quit"); err != nil {
		return
	}
	quitAction.Triggered().Attach(func() { dlg.Close() })
	if err = dlg.Menu().Actions().Add(quitAction); err != nil {
		return
	}

	// Connect button
	connectAction := walk.NewAction()
	if err = connectAction.SetText("&Connect"); err != nil {
		return
	}
	connectAction.Triggered().Attach(func() {
		result, addr, pw, err := runConnectDialog(addressStr, password, dlg)
		if err != nil {
			uiLogError(fmt.Sprintf("Failed to run connect dialog: %s", err))
			return
		}
		if result {
			if err = initSocketAddr(addr); err != nil {
				uiLogError(fmt.Sprintf("Couldn't use that address: %s", err))
				return
			}
			password = pw
			dlg.ui.rconOutput.SetText("")
		}
	})
	if err = dlg.Menu().Actions().Add(connectAction); err != nil {
		return
	}

	// Handle input
	dlg.ui.rconInput.KeyPress().Attach(func(key walk.Key) {
		if key != walk.KeyReturn {
			return
		}

		if address == nil {
			uiLogError("No server configured.")
			return
		}

		cmd := dlg.ui.rconInput.Text()
		dlg.ui.rconInput.SetText("")

		uiLog(address.String() + "> " + cmd)
		sendRcon(cmd)
	})

	// When window is initialized we can let a secondary routine print all
	// output received
	dlg.Synchronize(func() {
		go func() {
			for {
				msg, err := receiveRcon()
				if err != nil {
					uiLogError(err.Error())
					continue
				}
				switch strings.ToLower(msg.Name) {
				case "print":
					uiLog(string(msg.Data))
				default:
					log.Println(msg.Name)
				}
			}
		}()
	})

	// Get rid of the console window
	freeConsole.Call()

	dlg.Show()

	// Message loop starts here and will block the main goroutine!
	dlg.Run()

	return
}
